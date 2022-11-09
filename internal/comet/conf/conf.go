package conf

import (
	"time"
)

// Init init config.
func Init() (conf *Config, err error) {
	conf = Default()
	return
}

// Default new a config with specified defualt value.
func Default() *Config {
	return &Config{
		Network: &Network{
			Tcp: &TCP{
				Bind:    []string{":5311"},
				ReadBuf: 4096,
				SendBuf: 4096,
			},
			TimerSize:        32,
			HandshakeTimeout: time.Second * 8,
		},
		Bucket: &Bucket{
			Size:          32,
			Channel:       1024,
			Group:         1024,
			RoutineAmount: 32,
			RoutineSize:   1024,
		},
	}
}

// Config is comet config.
type Config struct {
	Network *Network
	Bucket  *Bucket
}

type Network struct {
	Tcp              *TCP
	HandshakeTimeout time.Duration
	TimerSize        int32
}

// TCP is tcp config.
type TCP struct {
	Bind    []string
	ReadBuf int
	SendBuf int
}

// Websocket is websocket config.
type Websocket struct {
	Bind []string
}

// Bucket is bucket config.
type Bucket struct {
	Size          int
	Channel       int
	Group         int
	RoutineAmount uint64
	RoutineSize   int
}
