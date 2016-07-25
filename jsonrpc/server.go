package jsonrpc

import (
	"encoding/json"
)

type message struct {
	Id *json.RawMessage `json:"id"`
}

type request struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type response struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}

func (r *router) sendResponses() {
	for {
		resp := <-r.down
		r.conn.WriteJSON(resp)
	}
}

func (r *router) Serve() {
	for {
		var raw json.RawMessage
		r.conn.ReadJSON(&raw)
		req, resp, err := guessRequestResponse(raw)
		if err != nil {
			r.error(err)
			continue
		}
		if req != nil {
			r.request <- req
		} else { // it's a response
			if resp == nil {
				panic("nil response")
			}
			r.response <- resp
		}
	}
}

func (r *router) error(err *jsonrpcerror) {
	r.down <- response{
		Id:    err.id,
		Error: err.message,
	}
}
