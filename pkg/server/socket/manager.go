package socket

import (
	"errors"
	"github.com/zhenjl/cityhash"
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/options"
)

type serverManager struct {
	buckets []*Bucket
	lis     Listener
	servers map[string]network.Server
	opt     *options.Options
}

func NewManager(lis Listener, opts ...options.Option) Manager {
	m := &serverManager{lis: lis}
	// set options
	opt := options.DefaultOptions()
	opt.ParseOptions(opts...)
	// channel buckets
	m.opt = opt
	m.buckets = make([]*Bucket, opt.BucketSize)
	for i := uint32(0); i < opt.BucketSize; i++ {
		m.buckets[i] = newBucket(opt.ChannelSize)
	}
	m.servers = make(map[string]network.Server)
	return m
}

func (m *serverManager) register(srv network.Server) {
	srvId := srv.Id()
	m.servers[srvId] = srv
	srv.SetHandler(m)
}

func (m *serverManager) Start() (err error) {
	if len(m.servers) == 0 {
		return errors.New("running server manager without servers")
	}
	for _, server := range m.servers {
		err = server.Start()
		if err != nil {
			return
		}
	}
	return
}

func (m *serverManager) Close() (err error) {
	for _, server := range m.servers {
		err = server.Close()
		if err != nil {
			return
		}
	}
	return
}

func (m *serverManager) Bucket(channelId string) *Bucket {
	idx := cityhash.CityHash32([]byte(channelId), uint32(len(channelId))) % m.opt.BucketSize
	return m.buckets[idx]
}

func (m *serverManager) Channel(channelId string) Channel {
	return m.Bucket(channelId).Channel(channelId)
}
