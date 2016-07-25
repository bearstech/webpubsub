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
			r.error(err.id, err.message)
			continue
		}
		if req != nil {
			r.request <- req
		} else { // it's a response
			r.response <- resp
		}
	}
}

func (r *router) error(id *json.RawMessage, msg string) {
	r.down <- response{
		Id:    id,
		Error: msg,
	}
}
