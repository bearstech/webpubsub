package jsonrpc

import (
	"encoding/json"
	"fmt"
)

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
			r.down <- err.Message()
			continue
		}
		if req != nil {
			r.request <- req
		} else { // it's a response
			if resp == nil {
				panic("nil response")
			}
			if resp.Id == nil {
				r.notification <- resp
			}
			var id uint64
			oups := json.Unmarshal(*resp.Id, &id)
			if oups != nil {
				r.down <- NewError(nil, oups.Error()).Message()
				continue
			}
			re, ok := r.response[id]
			if !ok {
				r.down <- NewError(nil, fmt.Sprintf("Unknwon id : %i", id))
				continue
			}
			re <- resp
		}
	}
}
