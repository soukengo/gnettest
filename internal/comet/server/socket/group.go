package socket

import (
	"errors"
	"gnettest/internal/comet/model"
	"sync"
)

var (
	ErrGroupDroped = errors.New("group droped")
)

type Group struct {
	ID       string
	rLock    sync.RWMutex
	channels map[string]*Channel
	drop     bool
	Online   int32 // dirty read is ok
}

func NewGroup(id string) (r *Group) {
	r = new(Group)
	r.ID = id
	r.drop = false
	r.Online = 0
	r.channels = make(map[string]*Channel)
	return
}

func (r *Group) Put(ch *Channel) (err error) {
	r.rLock.Lock()
	if !r.drop {
		if r.channels[ch.Uid] == nil {
			r.channels[ch.Uid] = ch
			r.Online++
		}
	} else {
		err = ErrGroupDroped
	}
	r.rLock.Unlock()
	return
}

func (r *Group) Del(ch *Channel) bool {
	r.rLock.Lock()
	delete(r.channels, ch.Uid)
	r.Online--
	r.drop = r.Online == 0
	r.rLock.Unlock()
	return r.drop
}

func (r *Group) Push(p *model.Packet) {
	r.rLock.RLock()
	for _, ch := range r.channels {
		_ = ch.Push(p)
	}
	r.rLock.RUnlock()
}

func (r *Group) Close() {
	r.rLock.RLock()
	for _, ch := range r.channels {
		ch.DelGroup(r.ID)
	}
	r.rLock.RUnlock()
}
