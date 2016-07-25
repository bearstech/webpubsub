package jsonrpc

import (
	"encoding/json"
	"testing"
)

func TestGuess(t *testing.T) {
	var raw json.RawMessage

	json.Unmarshal([]byte("{\"id\": 42, \"method\": \"plop\""), &raw)

	req, resp, err := guessRequestResponse(raw)
	if resp != nil {
		t.Error("resp")
	}
	if err == nil {
		t.Error("err")
	}
	if req != nil {
		t.Error("req")
	}

	json.Unmarshal([]byte("{\"id\": 42, \"method\": \"plop\"}"), &raw)
	req, resp, err = guessRequestResponse(raw)
	if resp != nil {
		t.Error("resp")
	}
	if err != nil {
		t.Error("err", err)
	}
	if req == nil {
		t.Error("req")
	}

	json.Unmarshal([]byte("{\"id\": 42, \"result\": \"plop\"}"), &raw)
	req, resp, err = guessRequestResponse(raw)
	if resp == nil {
		t.Error("resp")
	}
	if err != nil {
		t.Error("err", err)
	}
	if req != nil {
		t.Error("req")
	}
}
