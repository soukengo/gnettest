package socket

import (
	"github.com/google/uuid"
	"sync"
)

type Bucket struct {
	channelSize uint32
	id          string
	chs         map[string]Channel
	cLock       sync.RWMutex
}

func newBucket(channelSize uint32) *Bucket {
	b := &Bucket{
		id:  uuid.New().String(),
		chs: make(map[string]Channel, channelSize),
	}
	return b
}

func (b *Bucket) PutChannel(ch Channel) {
	b.cLock.Lock()
	b.chs[ch.Id()] = ch
	b.cLock.Unlock()
}
func (b *Bucket) DelChannel(channelId string) {
	b.cLock.Lock()
	delete(b.chs, channelId)
	b.cLock.Unlock()
}

func (b *Bucket) Channel(channelId string) (ch Channel) {
	b.cLock.RLock()
	ch = b.chs[channelId]
	b.cLock.RUnlock()
	return
}
