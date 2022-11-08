package nbio

import (
	"bufio"
	"bytes"
	"github.com/google/uuid"
	"github.com/lesismal/nbio"
	"gnettest/pkg/server/socket/packet"
	"net"
	"sync/atomic"
)

type nbioConn struct {
	id     string
	core   *nbio.Conn
	parser *packet.Parser
	buf    *bytes.Buffer
	reader *bufio.Reader
	closed int32
}

func newConn(gc *nbio.Conn, parser *packet.Parser) *nbioConn {
	buf := bytes.NewBuffer([]byte{})
	return &nbioConn{
		id:     uuid.New().String(),
		core:   gc,
		parser: parser,
		buf:    buf,
		reader: bufio.NewReader(buf),
		closed: 0,
	}
}

func (g *nbioConn) Id() string {
	return g.id
}

func (g *nbioConn) Read() (p packet.IPacket, err error) {
	return g.parser.Parse(g.Id(), g.reader)
}

func (g *nbioConn) Send(packet packet.IPacket) (err error) {
	err = packet.PackTo(g)
	return
}

func (g *nbioConn) Write(data []byte) (n int, err error) {
	return g.core.Write(data)
}

func (g *nbioConn) RemoteAddr() net.Addr {
	return g.core.RemoteAddr()
}
func (g *nbioConn) Close() error {
	if atomic.LoadInt32(&g.closed) > 0 {
		return nil
	}
	if !g.markClosed() {
		return nil
	}
	return g.core.Close()
}

func (g *nbioConn) IsAsync() bool {
	return false
}

func (g *nbioConn) markClosed() bool {
	return atomic.CompareAndSwapInt32(&g.closed, 0, 1)
}
