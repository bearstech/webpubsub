package jsonrpc

import (
	"encoding/json"
)

type Client struct {
	router *router
}

func (r *router) Client() *Client {
	return &Client{r}
}

func (c *Client) Notification(serviceMethod string, args interface{}) error {
	return c.call(0, serviceMethod, args, nil)
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	c.router.down_mutex.Lock()
	c.router.down_id++
	id := c.router.down_id
	c.router.down_mutex.Unlock()
	return c.call(id, serviceMethod, args, reply)
}

func (c *Client) call(id uint64, serviceMethod string, args interface{}, reply interface{}) error {

	p, err := json.Marshal(args)
	if err != nil {
		return nil
	}
	var i []byte
	if id == 0 {
		i, err = json.Marshal(nil)
	} else {
		i, err = json.Marshal(id)
	}
	if err != nil {
		return err
	}
	pp := json.RawMessage(p)
	ii := json.RawMessage(i)

	req := request{
		Method: serviceMethod,
		Params: &pp,
		Id:     &ii,
	}
	c.router.down <- req
	// FIXME wait for the response
	return nil
}
