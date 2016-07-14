package mailbox

import (
	"testing"
	"time"
)

func TestMailbox(t *testing.T) {
	m := New(time.Duration(150*time.Millisecond), 100)
	cpt := 0
	machin := make(chan bool)
	mxAlice := m.Subscribe("alice")

	go func() {
		for {
			msg := <-mxAlice.Mails()
			if string(msg) == "plop" {
				cpt++
				machin <- true
			}
		}
	}()
	mxBob := m.Subscribe("bob")
	go func() {
		for {
			msg := <-mxBob.Mails()
			if string(msg) == "plop" {
				cpt++
				machin <- true
			}
		}
	}()
	published := m.Publish([]byte("plop"))
	if published != 2 {
		t.Error("Message sent", published)
	}
	<-machin
	<-machin

	if cpt != 2 {
		t.Error("Not enough messages sent", cpt)
	}

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
	mxAlice2 := m.Subscribe("alice")
	mxAlice2.Mails() <- []byte("plop")
	<-machin
	if cpt != 3 {
		t.Error("Direct mail miss : ", cpt)
	}

	<-m.dead
	if m.Length() > 0 {
		t.Error("There are zombies")
	}

}
