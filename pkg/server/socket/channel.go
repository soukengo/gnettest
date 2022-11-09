package socket

import (
	"context"
	log "github.com/golang/glog"
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/packet"
	"net"
	"sync"
)

type channel struct {
	clientIP  string
	conn      network.Connection
	done      *sync.Once
	recvQueue chan packet.IPacket
	sendQueue chan packet.IPacket
	ctx       context.Context
	cancel    context.CancelFunc
	received  func(p packet.IPacket)
}

func newChannel(conn network.Connection, recvQueueSize uint32, sendQueueSize uint32, received func(p packet.IPacket)) *channel {
	ctx, cancel := context.WithCancel(context.Background())
	ins := &channel{
		ctx:       ctx,
		cancel:    cancel,
		done:      new(sync.Once),
		recvQueue: make(chan packet.IPacket, recvQueueSize),
		sendQueue: make(chan packet.IPacket, sendQueueSize),
		conn:      conn,
		received:  received,
	}
	ins.clientIP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	go ins.recvLoop()
	go ins.dispatch()
	return ins
}

func (c *channel) Id() string {
	return c.conn.Id()
}

func (c *channel) ClientIP() string {
	return c.clientIP
}

func (c *channel) Send(data packet.IPacket) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		c.sendQueue <- data
	}
	return nil
}

func (c *channel) Receive(p packet.IPacket) {
	select {
	case <-c.ctx.Done():
		return
	default:
		c.recvQueue <- p
	}
}

func (c *channel) Close() (err error) {
	c.done.Do(func() {
		c.cancel()
		err = c.conn.Close()
	})
	return
}

// recvLoop read packet from connection
func (c *channel) recvLoop() {
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

// dispatch send packet to connection
func (c *channel) dispatch() {
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
				log.Errorf("conn.Send err: %v", err)
				return
			}
		}
	}
}
