package mailbox

import (
	"sync"
	"testing"
	"time"
)

func TestMailbox(t *testing.T) {
	m := New(time.Duration(150*time.Millisecond), 100)
	var wg sync.WaitGroup
	mxAlice := m.Subscribe("alice")

	wg.Add(2)
	go func() {
		mails := mxAlice.Mails()
		for {
			msg := <-mails
			if string(msg) == "plop" {
				wg.Done()
			}
		}
	}()
	mxBob := m.Subscribe("bob")
	go func() {
		mails := mxBob.Mails()
		for {
			msg := <-mails
			if string(msg) == "plop" {
				wg.Done()
			}
		}
	}()
	published := m.Publish("a", []byte("plop"))
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
	mxAlice2.Mails() <- []byte("plop")
	wg.Wait()

	<-m.dead
	if m.Length() > 0 {
		t.Error("There are zombies")
	}

}
