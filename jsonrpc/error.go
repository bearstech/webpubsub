package jsonrpc

import (
	"encoding/json"
)

type jsonrpcerror struct {
	id      *json.RawMessage
	message string
}

func (e *jsonrpcerror) Error() string {
	return e.message
}

func (e *jsonrpcerror) Message() *response {
	return &response{
		Id:    e.id,
		Error: e.message,
	}
}

func NewError(id *json.RawMessage, message string) *jsonrpcerror {
	return &jsonrpcerror{id, message}
}
