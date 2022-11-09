package socket

import (
	"gnettest/internal/comet/model"
	"gnettest/pkg/server/socket"
	"sync"
	"time"
)

type Channel struct {
	bucket    *Bucket
	groups    map[string]*Group
	LastHb    time.Time
	core      socket.Channel
	mutex     sync.Mutex
	timer     timer
	active    time.Duration
	Uid       string
	Key       string
	groupLock sync.Mutex
	closeOnce *sync.Once
}

// timer 计时器接口
type timer interface {
	Reset(d time.Duration) bool
	Stop() bool
}

// NewChannel new a channel.
func NewChannel(core socket.Channel) *Channel {
	if core == nil {
		return nil
	}
	c := new(Channel)

	c.groups = make(map[string]*Group)
	c.core = core
	c.Key = core.Id()
	c.closeOnce = new(sync.Once)
	return c
}

func (c *Channel) ClientIp() string {
	return c.core.ClientIP()
}

// Enable 使用原生的计时器
func (c *Channel) Enable(alive time.Duration, onExpired func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.active = alive
	c.LastHb = time.Now()
	c.timer = time.AfterFunc(c.active, onExpired)
}

func (c *Channel) JoinGroup(ids ...string) {
	for _, rid := range ids {
		_ = c.bucket.JoinGroup(rid, c)
	}
}

func (c *Channel) QuitGroup(ids ...string) {
	for _, rid := range ids {
		c.bucket.QuitGroup(rid, c)
	}
}

// Renew 续期
func (c *Channel) Renew() {
	c.timer.Reset(c.active)
	c.LastHb = time.Now()
}

// ResetTimer 重置timer
func (c *Channel) ResetTimer(duration time.Duration) {
	c.timer.Reset(duration)
}

func (c *Channel) Push(p *model.Packet) (err error) {
	return c.core.Send(p)
}

func (c *Channel) Close() (err error) {
	c.closeOnce.Do(func() {
		c.timer.Stop()
		err = c.core.Close()
	})
	return
}

func (c *Channel) AddGroup(group *Group) {
	c.groupLock.Lock()
	c.groups[group.ID] = group
	c.groupLock.Unlock()
}

func (c *Channel) DelGroup(rid string) {
	c.groupLock.Lock()
	delete(c.groups, rid)
	c.groupLock.Unlock()
}
