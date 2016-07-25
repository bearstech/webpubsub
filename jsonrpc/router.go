package jsonrpc

import (
	"encoding/json"
	"sync"
)

type JSONReaderWriter interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
}

type router struct {
	conn         JSONReaderWriter
	id           uint64
	request      chan *request
	response     map[uint64]chan *response
	notification chan *response
	down         chan interface{}
	down_mutex   sync.Mutex
	down_id      uint64
}

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
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
}

func NewRouter(rw JSONReaderWriter) *router {
	r := router{
		conn:         rw,
		request:      make(chan *request),
		notification: make(chan *response),
		down:         make(chan interface{}),
		response:     make(map[uint64]chan *response),
	}
	go r.sendResponses()
	return &r
}

func guessRequestResponse(raw json.RawMessage) (req *request, resp *response, jsonerr *jsonrpcerror) {
	var stuff map[string]interface{}
	err := json.Unmarshal(raw, &stuff)
	if err != nil {
		return nil, nil, NewError(nil, err.Error())
	}
	_, ok := stuff["method"]
	if ok { // it should be a request
		err = json.Unmarshal(raw, &req)
		if err != nil {
			return nil, nil, NewError(nil, err.Error())
		}
		return req, nil, nil
	}
	_, err_ok := stuff["error"]
	_, result_ok := stuff["result"]
	var b message
	json.Unmarshal(raw, &b)
	id := b.Id
	if !err_ok && !result_ok {
		// impossible
		return nil, nil, NewError(id, "Both result and error")
	}
	if !(err_ok || result_ok) {
		return nil, nil, NewError(id, "Can't find result nor error")
	}
	// it's a response
	json.Unmarshal(raw, &resp)
	return
}
