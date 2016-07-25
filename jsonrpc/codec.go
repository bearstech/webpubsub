package jsonrpc

import (
	"encoding/json"
	"errors"
	"net/rpc"
	"sync"
)

var errMissingParams = errors.New("jsonrpc: request body missing params")

type serverCodec struct {
	req     *request
	resp    *response
	route   *router
	mutex   sync.Mutex
	seq     uint64
	pending map[uint64]*json.RawMessage
}

func NewServerCodec(route *router) *serverCodec {
	c := serverCodec{
		route:   route,
		pending: make(map[uint64]*json.RawMessage),
	}
	return &c
}

func (c *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	c.req = <-c.route.request
	r.ServiceMethod = c.req.Method

	c.mutex.Lock()
	c.seq++
	c.pending[c.seq] = c.req.Id
	c.req.Id = nil
	r.Seq = c.seq
	c.mutex.Unlock()
	return nil
}

func (c *serverCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	if c.req.Params == nil {
		return errMissingParams
	}
	var params [1]interface{}
	params[0] = x
	return json.Unmarshal(*c.req.Params, &params)
}

var null = json.RawMessage([]byte("null"))

func (c *serverCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	c.mutex.Lock()
	b, ok := c.pending[r.Seq]
	if !ok {
		c.mutex.Unlock()
		return errors.New("invalid sequence number in response")
	}
	delete(c.pending, r.Seq)
	c.mutex.Unlock()
	if b == nil {
		b = &null
	}

	c.resp = &response{Id: b}
	if r.Error == "" {

		raw, err := json.Marshal(x)
		if err == nil {
			rr := json.RawMessage(raw)
			c.resp.Result = &rr
		} else {
			return err
		}
	} else {
		c.resp.Error = r.Error
	}
	c.route.down <- c.resp
	return nil
}

func (c *serverCodec) Close() error {
	return nil
}
