package mailbox

import (
	"sync"
	"time"
)

var zero time.Time

func init() {
	zero = time.Unix(0, 0)
}

type Mailbox struct {
	eta       time.Time
	ttl       time.Duration
	trashtube chan *Mailbox
	user      string
	Mails     chan []byte
}

func (mx *Mailbox) Leave() {
	mx.eta = time.Now().Add(mx.ttl)
	go func() {
		time.Sleep(mx.ttl)
		mx.trashtube <- mx
	}()
}

type Mailboxes struct {
	boxes     map[string]*Mailbox
	lock      sync.Mutex
	ttl       time.Duration
	boxSize   int
	trashtube chan *Mailbox
}

func New(ttl time.Duration, boxSize int) *Mailboxes {
	m := Mailboxes{
		boxes:     make(map[string]*Mailbox),
		ttl:       ttl,
		boxSize:   boxSize,
		trashtube: make(chan *Mailbox),
	}
	go m.gc()
	return &m
}

func (m *Mailboxes) Length() int {
	return len(m.boxes)
}

func (m *Mailboxes) gc() {
	for {
		mx := <-m.trashtube
		now := time.Now()
		if mx.eta != zero && now.After(mx.eta) {
			defer m.lock.Unlock()
			m.lock.Lock()
			delete(m.boxes, mx.user)
		}
	}
}

func (m *Mailboxes) Publish(mail []byte) {
	defer m.lock.Unlock()
	m.lock.Lock()
	for _, box := range m.boxes {
		box.Mails <- mail
	}
}

func (m *Mailboxes) Subscribe(user string) *Mailbox {
	defer m.lock.Unlock()
	m.lock.Lock()
	var mx *Mailbox
	mx, ok := m.boxes[user]
	if !ok {
		mx = &Mailbox{
			Mails:     make(chan []byte, m.boxSize),
			trashtube: m.trashtube,
			user:      user,
			eta:       zero,
		}
		m.boxes[user] = mx
	} else {
		mx.eta = zero
	}
	return mx
}
