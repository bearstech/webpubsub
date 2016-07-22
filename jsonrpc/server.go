package jsonrpc

import (
	"encoding/json"
)

type JSONReaderWriter interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
}

type router struct {
	conn     JSONReaderWriter
	request  chan request
	response chan response
}

type request struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     *json.RawMessage `json:"id"`
}

type response struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}

type berk struct {
	Id *json.RawMessage `json:"id"`
}

func NewRouter(rw JSONReaderWriter) *router {
	r := router{
		conn:     rw,
		request:  make(chan request),
		response: make(chan response),
	}
	return &r
}

func (r *router) sendResponses() {
	for {
		resp := <-r.response
		r.conn.WriteJSON(resp)
	}
}

func (r *router) Serve() {
	go r.sendResponses()
	for {
		var raw json.RawMessage
		var stuff map[string]interface{}
		r.conn.ReadJSON(&raw)

		err := json.Unmarshal(raw, &stuff)
		if err != nil {
			r.error(nil, "Can't parse JSON : "+err.Error())
			continue
		}
		_, ok := stuff["method"]
		if ok { // it should be a request
			var req request
			json.Unmarshal(raw, &req)
			r.request <- req
		} else {
			_, err_ok := stuff["error"]
			_, result_ok := stuff["result"]
			var b berk
			json.Unmarshal(raw, &b)
			id := b.Id
			if err_ok && result_ok {
				// impossible
				r.error(id, "Both result and error")
				continue
			}
			if !(err_ok || result_ok) {
				r.error(id, "Can't find result nor error")
				continue
			}
			// it's a response
			var resp response
			json.Unmarshal(raw, &resp)
			r.response <- resp
		}
	}
}

func (r *router) error(id *json.RawMessage, msg string) {
	r.response <- response{
		Id:    id,
		Error: msg,
	}
}
