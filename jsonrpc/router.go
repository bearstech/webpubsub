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
	conn       JSONReaderWriter
	id         uint64
	request    chan *request
	response   chan *response
	down       chan interface{}
	down_mutex sync.Mutex
	down_id    uint64
}

func NewRouter(rw JSONReaderWriter) *router {
	r := router{
		conn:     rw,
		request:  make(chan *request),
		response: make(chan *response),
		down:     make(chan interface{}),
	}
	go r.sendResponses()
	return &r
}

type jsonrpcerror struct {
	id      *json.RawMessage
	message string
}

func (e *jsonrpcerror) Error() string {
	return e.message
}

func guessRequestResponse(raw json.RawMessage) (req *request, resp *response, jsonerr *jsonrpcerror) {
	var stuff map[string]interface{}
	err := json.Unmarshal(raw, &stuff)
	if err != nil {
		return nil, nil, &jsonrpcerror{nil, err.Error()}
	}
	_, ok := stuff["method"]
	if ok { // it should be a request
		err = json.Unmarshal(raw, &req)
		if err != nil {
			return nil, nil, &jsonrpcerror{nil, err.Error()}
		}
		return req, nil, nil
	}
	_, err_ok := stuff["error"]
	_, result_ok := stuff["result"]
	var b message
	json.Unmarshal(raw, &b)
	id := b.Id
	if err_ok && result_ok {
		// impossible
		return nil, nil, &jsonrpcerror{id, "Both result and error"}
	}
	if !(err_ok || result_ok) {
		return nil, nil, &jsonrpcerror{id, "Can't find result nor error"}
	}
	// it's a response
	json.Unmarshal(raw, &resp)
	return
}
