package mailbox

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestEncoding(t *testing.T) {
	type test struct {
		Age  int    `json:age`
		Name string `json:name`
	}
	msg := Message{
		Path: "a",
		Body: test{
			Age:  42,
			Name: "Robert",
		},
	}
	j, _ := json.Marshal(msg)
	fmt.Println(string(j))

}

func TestMailbox(t *testing.T) {
	m := New(time.Duration(150*time.Millisecond), 100)
	var wg sync.WaitGroup
	mxAlice := m.Subscribe("alice")

	wg.Add(2)
	go func() {
		mails := mxAlice.Mails()
		for {
			msg := <-mails
			if msg.Body == "plop" {
				wg.Done()
			}
		}
	}()
	mxBob := m.Subscribe("bob")
	go func() {
		mails := mxBob.Mails()
		for {
			msg := <-mails
			if msg.Body == "plop" {
				wg.Done()
			}
		}
	}()
	published := m.Publish(Message{"a", "plop"})
	if published != 2 {
		t.Error("Message sent", published)
	}
	wg.Wait()

	mxBob.Leave()

	resp := <-m.dead
	if resp != "bob" {
		t.Error("Bob should be deleted : ", resp)
	}
	if m.Length() != 1 {
		t.Error("One should had left", m.Length())
	}
	if !m.boxes["alice"].eta.IsZero() {
		t.Error("Alice ETA is not zero")
	}
	mxAlice.Leave()
	if m.boxes["alice"].eta.IsZero() {
		t.Error("Alice ETA is zero")
	}
	mxAlice.DontLeave()
	if !m.boxes["alice"].eta.IsZero() {
		t.Error("Alice ETA is not zero")
	}

	mxAlice.Leave()
	if mxAlice.ETA().IsZero() {
		t.Error("Alice ETA is zero")
	}
	mxAlice2 := m.Subscribe("alice")
	if !mxAlice.ETA().IsZero() {
		t.Error("Alice ETA is not zero")
	}
	wg.Add(1)
	mxAlice2.Mails() <- Message{"a", "plop"}
	wg.Wait()

	<-m.dead
	if m.Length() > 0 {
		t.Error("There are zombies")
	}

}
