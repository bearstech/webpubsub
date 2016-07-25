package jsonrpc

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

type goWsRw struct {
	conn *websocket.Conn
}

func (rw *goWsRw) ReadJSON(v interface{}) error {
	return websocket.JSON.Receive(rw.conn, v)
}

func (rw *goWsRw) WriteJSON(v interface{}) error {
	return websocket.JSON.Send(rw.conn, v)
}

type Args struct {
	A int
	B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func simpleSever(ws *websocket.Conn) {
	r := NewRouter(&goWsRw{ws})
	c := NewServerCodec(r)
	s := rpc.DefaultServer
	go s.ServeCodec(c)
	r.Serve()
}

func TestSimpleServer(t *testing.T) {
	go func() {
		rpc.Register(new(Arith))
		http.Handle("/conn", websocket.Handler(simpleSever))
		http.ListenAndServe("localhost:7000", nil)
	}()
	origin := "http://localhost/"
	url := "ws://localhost:7000/conn"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}
	client := jsonrpc.NewClient(ws)
	var result int
	client.Call("Arith.Multiply", Args{A: 4, B: 5}, &result)
	if result != 20 {
		t.Error("Bad response :", result)
	}
}

func talkingToClient(ws *websocket.Conn) {
	r := NewRouter(&goWsRw{ws})
	c := r.Client()
	c.Notification("chan", "beuha")
}

func TestServerToClientNotification(t *testing.T) {
	go func() {
		http.Handle("/conn2", websocket.Handler(talkingToClient))
		http.ListenAndServe("localhost:7000", nil)
	}()

	origin := "http://localhost/"
	url := "ws://localhost:7000/conn2"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}
	var notif request
	websocket.JSON.Receive(ws, &notif)
	if notif.Method != "chan" {
		t.Fatal(notif)
	}
}

func askSomethingToClient(ws *websocket.Conn) {
	r := NewRouter(&goWsRw{ws})
	c := r.Client()
	var resp string
	err := c.Call("hello", "world", &resp)
	if err != nil {
		panic(err)
	}
	if resp != "beuha" {
		panic("not Beuha")
	}
	c.Notification("chan", true)
}

func TestServerToClientRequest(t *testing.T) {
	go func() {
		http.Handle("/conn3", websocket.Handler(askSomethingToClient))
		http.ListenAndServe("localhost:7000", nil)
	}()

	origin := "http://localhost/"
	url := "ws://localhost:7000/conn2"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}
	var req request
	var nianiani = "beuha"
	websocket.JSON.Receive(ws, &req)
	j, err := json.Marshal(nianiani)
	jj := json.RawMessage(j)
	resp := response{
		Id:     req.Id,
		Result: &jj,
	}
	websocket.JSON.Send(ws, resp)
	var ok bool
	websocket.JSON.Receive(ws, &ok)
	if !ok {
		t.Error("oups")
	}
}
