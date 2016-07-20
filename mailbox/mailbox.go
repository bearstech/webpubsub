package mailbox

import (
	"sync"
	"time"
)

type Matcher interface {
	Match(path string) bool
}

type AllMatcher struct {
}

func (am *AllMatcher) Match(path string) bool {
	return true
}

type Message struct {
	Path string
	Body interface{}
}

type user string

type Mailboxes struct {
	boxes   map[user]*mailbox
	lock    sync.RWMutex
	ttl     time.Duration
	boxSize int
	dead    chan user
}

type mailbox struct {
	death   *time.Timer
	eta     time.Time
	mails   chan Message
	pattern Matcher
}

type MailboxProxy struct {
	parent *Mailboxes
	user   user
}

func (mp *MailboxProxy) Mails() chan Message {
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
		boxes:   make(map[user]*mailbox),
		ttl:     ttl,
		boxSize: boxSize,
		dead:    make(chan user),
	}
	return &m
}

func (m *Mailboxes) Length() int {
	return len(m.boxes)
}

func (m *Mailboxes) Publish(mail Message) int {
	defer m.lock.RUnlock()
	m.lock.RLock()
	cpt := 0
	for _, box := range m.boxes {
		if box.pattern.Match(mail.Path) {
			box.mails <- mail
			cpt += 1
		}
	}
	return cpt
}

func (m *Mailboxes) Subscribe(user string) *MailboxProxy {
	return m.SubscribePattern(user, &AllMatcher{})
}

func (m *Mailboxes) SubscribePattern(uzer string, pattern Matcher) *MailboxProxy {
	m.lock.RLock()
	mx, ok := m.boxes[user(uzer)]
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
		m.boxes[user(uzer)] = &mailbox{
			mails:   make(chan Message, m.boxSize),
			pattern: pattern,
		}
		m.lock.Unlock()
	}
	return &MailboxProxy{
		user:   user(uzer),
		parent: m,
	}
}
