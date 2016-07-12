package mailbox

import (
	"fmt"
	"testing"
	"time"
)

func TestMailbox(t *testing.T) {
	m := New(time.Duration(10*time.Second), 100)
	mx := m.Subscribe("box")
	go func() {
		for {
			msg := <-mx.Mails
			fmt.Println(string(msg))
		}
	}()
	m.Publish([]byte("plop"))
	mx.Leave()
	fmt.Println(m.boxes)
	time.Sleep(time.Duration(11 * time.Second))
	fmt.Println(m.boxes)
}
