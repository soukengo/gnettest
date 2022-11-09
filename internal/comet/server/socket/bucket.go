package socket

import (
	"gnettest/internal/comet/model"
	"sync"
	"sync/atomic"

	"gnettest/internal/comet/conf"
)

type Bucket struct {
	c           *conf.Bucket
	cLock       sync.RWMutex
	chs         map[string]*Channel
	groups      map[string]*Group
	routines    []chan *model.PushGroupReq
	routinesNum uint64
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, c.Channel)
	b.c = c
	b.groups = make(map[string]*Group, c.Group)
	b.routines = make([]chan *model.PushGroupReq, c.RoutineAmount)
	for i := uint64(0); i < c.RoutineAmount; i++ {
		c := make(chan *model.PushGroupReq, c.RoutineSize)
		b.routines[i] = c
		go b.groupProc(c)
	}
	return
}

// ChannelCount channel count in the bucket
func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

// Put put a channel according with sub key.
func (b *Bucket) Put(ch *Channel) (err error) {
	var (
		group *Group
	)
	ch.bucket = b
	b.cLock.Lock()
	b.chs[ch.Key] = ch
	b.cLock.Unlock()
	if group != nil {
		err = group.Put(ch)
	}
	return
}

func (b *Bucket) JoinGroup(rid string, ch *Channel) (err error) {
	var (
		group *Group
		ok    bool
	)
	b.cLock.Lock()
	if group, ok = b.groups[rid]; !ok {
		group = NewGroup(rid)
		b.groups[rid] = group
	}
	ch.AddGroup(group)
	b.cLock.Unlock()
	if group != nil {
		err = group.Put(ch)
	}
	return
}

func (b *Bucket) QuitGroup(rid string, ch *Channel) {
	b.cLock.Lock()
	group, ok := b.groups[rid]
	b.cLock.Unlock()
	if ok {
		if group.Del(ch) {
			b.DelGroup(group)
		}
	}
	ch.DelGroup(rid)
	return
}

func (b *Bucket) Del(dch *Channel) {
	var (
		ok     bool
		ch     *Channel
		groups map[string]*Group
	)
	b.cLock.Lock()
	if ch, ok = b.chs[dch.Key]; ok {
		groups = make(map[string]*Group)
		for id, r := range groups {
			groups[id] = r
		}
		if ch == dch {
			delete(b.chs, ch.Key)
		}
	}
	b.cLock.Unlock()
	if len(groups) > 0 {
		for _, group := range groups {
			if group.Del(ch) {
				b.DelGroup(group)
			}
		}
	}
}

func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) Broadcast(p *model.Packet, op int32) {
	var ch *Channel
	b.cLock.RLock()
	for _, ch = range b.chs {
		_ = ch.Push(p)
	}
	b.cLock.RUnlock()
}

func (b *Bucket) Group(rid string) (group *Group) {
	b.cLock.RLock()
	group = b.groups[rid]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) DelGroup(group *Group) {
	b.cLock.Lock()
	delete(b.groups, group.ID)
	b.cLock.Unlock()
	group.Close()
}

func (b *Bucket) PushGroup(arg *model.PushGroupReq) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	b.routines[num] <- arg
}

// groupProc
func (b *Bucket) groupProc(c chan *model.PushGroupReq) {
	for {
		arg := <-c
		if group := b.Group(arg.GroupID); group != nil {
			group.Push(arg.Packet)
		}
	}
}
