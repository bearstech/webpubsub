package mailbox

import (
	"testing"
	"time"
)

func TestMailbox(t *testing.T) {
	m := New(time.Duration(100*time.Millisecond), 100)
	mx_alice := m.Subscribe("alice")
	mx_bob := m.Subscribe("bob")
	cpt := 0
	go func() {
		for {
			msg := <-mx_alice.Mails
			if string(msg) != "plop" {
				t.Error("Bad message :", msg)
			}
			cpt += 1
		}
	}()
	go func() {
		for {
			<-mx_bob.Mails
			cpt += 1
		}
	}()
	m.Publish([]byte("plop"))
	time.Sleep(time.Duration(10 * time.Millisecond))
	if cpt != 2 {
		t.Error("Not enough messages sent", cpt)
	}
	mx_alice.Leave()
	mx_alice2 := m.Subscribe("alice")
	if mx_alice != mx_alice2 {
		t.Error("It's not twice the same box")
	}
	time.Sleep(time.Duration(200 * time.Millisecond))
	if m.Length() != 2 {
		t.Error("Cancel leave miss")
	}
	mx_bob.Leave()
	time.Sleep(time.Duration(200 * time.Millisecond))
	if m.Length() != 1 {
		t.Error("One should had left")
	}
}
