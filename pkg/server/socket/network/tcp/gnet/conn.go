package gnet

import (
	"bufio"
	"fmt"
	log "github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/panjf2000/gnet/v2"
	"gnettest/pkg/server/socket/packet"
	"go.uber.org/atomic"
	"io"
	"net"
	"sync"
)

var (
	maxWriteBufferSize = 1024 * 1024 * 1 // 1M
)

type tcpConn struct {
	id     string
	core   gnet.Conn
	parser *packet.Parser
	reader io.Reader
	closed *atomic.Bool
}

func newConn(gc gnet.Conn, parser *packet.Parser) *tcpConn {
	return &tcpConn{
		id:   uuid.New().String(),
		core: gc, parser: parser,
		reader: bufio.NewReader(gc),
		closed: atomic.NewBool(false),
	}
}

func (g *tcpConn) Id() string {
	return g.id
}

func (g *tcpConn) Read() (p packet.IPacket, err error) {
	return g.parser.Parse(g.Id(), g.reader)
}

func (g *tcpConn) Send(packet packet.IPacket) (err error) {
	err = packet.PackTo(g)
	return
}

func (g *tcpConn) Write(data []byte) (n int, err error) {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	err = g.core.AsyncWrite(data, func(c gnet.Conn) error {
		if c.OutboundBuffered() > maxWriteBufferSize {
			log.Errorf("OutboundBuffered size exceeds max write buffer size, current: %v", c.OutboundBuffered())
			c.Close()
		}
		wg.Done()
		return nil
	})
	wg.Wait()
	if err == nil {
		n = len(data)
	}
	return
}

func (g *tcpConn) RemoteAddr() net.Addr {
	return g.core.RemoteAddr()
}
func (g *tcpConn) Close() error {
	if g.closed.Load() {
		return nil
	}
	if !g.closed.CAS(false, true) {
		return nil
	}
	fmt.Println("close client")
	return g.core.Close()
}

func (g *tcpConn) IsAsync() bool {
	return false
}
