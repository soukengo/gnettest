package main

import (
	"context"
	log "github.com/golang/glog"
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/packet"
	"net"
	"sync"
)

type Channel struct {
	clientIP  string
	conn      network.Connection
	done      *sync.Once
	recvQueue chan packet.IPacket
	sendQueue chan packet.IPacket
	ctx       context.Context
	cancel    context.CancelFunc
	attrs     map[string]any
	received  func(p packet.IPacket)
}

func newChannel(conn network.Connection, recvQueueSize uint32, sendQueueSize uint32, received func(p packet.IPacket)) *Channel {
	ctx, cancel := context.WithCancel(context.Background())
	ins := &Channel{
		ctx:       ctx,
		cancel:    cancel,
		done:      new(sync.Once),
		recvQueue: make(chan packet.IPacket, recvQueueSize),
		sendQueue: make(chan packet.IPacket, sendQueueSize),
		conn:      conn,
		attrs:     map[string]any{},
		received:  received,
	}
	ins.clientIP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	go ins.recvLoop()
	go ins.dispatch()
	return ins
}

func (c *Channel) Id() string {
	return c.conn.Id()
}

func (c *Channel) ClientIP() string {
	return c.clientIP
}

func (c *Channel) Send(data packet.IPacket) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		c.sendQueue <- data
	}
	return nil
}

func (c *Channel) Receive(p packet.IPacket) {
	select {
	case <-c.ctx.Done():
		return
	default:
		c.recvQueue <- p
	}
}

func (c *Channel) Close() (err error) {
	c.done.Do(func() {
		c.cancel()
		err = c.conn.Close()
	})
	return
}

func (c *Channel) Attrs() map[string]any {
	return c.attrs
}

func (c *Channel) SetAttr(key string, value any) {
	c.attrs[key] = value
}

func (c *Channel) recvLoop() {
	defer func() {
		_ = c.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case p, ok := <-c.recvQueue:
			if !ok {
				return
			}
			c.received(p)
		}
	}
}

func (c *Channel) dispatch() {
	defer func() {
		_ = c.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case p := <-c.sendQueue:
			err := c.conn.Send(p)
			if err != nil {
				log.Errorf("Write err: %v", err)
				return
			}
		}
	}
}
