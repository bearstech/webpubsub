package mailbox

import (
	"sync"
	"time"
)

type Mailboxes struct {
	boxes   map[string]*mailbox
	lock    sync.RWMutex
	ttl     time.Duration
	boxSize int
	dead    chan string
}

type mailbox struct {
	death *time.Timer
	eta   time.Time
	mails chan []byte
}

type MailboxProxy struct {
	parent *Mailboxes
	user   string
}

func (mp *MailboxProxy) Mails() chan []byte {
	defer mp.parent.lock.RUnlock()
	mp.parent.lock.RLock()
	return mp.parent.boxes[mp.user].mails
}

func (mp *MailboxProxy) ETA() time.Time {
	defer mp.parent.lock.RUnlock()
	mp.parent.lock.RLock()
	return mp.parent.boxes[mp.user].eta
}

func (mp *MailboxProxy) Leave() {
	mp.parent.lock.RLock()
	m := mp.parent.boxes[mp.user]
	mp.parent.lock.RUnlock()
	now := time.Now()
	if m.death != nil {
		m.death.Reset(mp.parent.ttl)
	} else {
		m.death = time.AfterFunc(mp.parent.ttl, func() {
			mp.parent.lock.Lock()
			delete(mp.parent.boxes, mp.user)
			mp.parent.lock.Unlock()
			mp.parent.dead <- mp.user
		})
	}
	m.eta = now.Add(mp.parent.ttl)
}

func (mp *MailboxProxy) DontLeave() {
	mp.parent.lock.RLock()
	mp.parent.boxes[mp.user].death.Stop()
	mp.parent.boxes[mp.user].eta = time.Time{}
	mp.parent.lock.RUnlock()
}

func New(ttl time.Duration, boxSize int) *Mailboxes {
	m := Mailboxes{
		boxes:   make(map[string]*mailbox),
		ttl:     ttl,
		boxSize: boxSize,
		dead:    make(chan string),
	}
	return &m
}

func (m *Mailboxes) Length() int {
	return len(m.boxes)
}

func (m *Mailboxes) Publish(mail []byte) int {
	defer m.lock.RUnlock()
	m.lock.RLock()
	cpt := 0
	for _, box := range m.boxes {
		box.mails <- mail
		cpt += 1
	}
	return cpt
}

func (m *Mailboxes) Subscribe(user string) *MailboxProxy {
	m.lock.RLock()
	mx, ok := m.boxes[user]
	if ok {
		gni := mx.death.Stop()
		if gni {
			mx.eta = time.Time{}
			// Start ugly hack
			// There is a deadlock without that
			gni = mx.death.Reset(time.Millisecond)
			if gni {
				panic("This timer should be stopped")
			}
			// End ugly hack
		}
	}
	m.lock.RUnlock()
	if !ok {
		m.lock.Lock()
		m.boxes[user] = &mailbox{
			mails: make(chan []byte, m.boxSize),
		}
		m.lock.Unlock()
	}
	return &MailboxProxy{
		user:   user,
		parent: m,
	}
}
