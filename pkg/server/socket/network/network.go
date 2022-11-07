package network

import (
	"gnettest/pkg/server/socket/packet"
	"net"
	"net/url"
)

type Server interface {
	Id() string
	Start() error
	Close() error
	EndPoint() *url.URL
	SetHandler(listener Handler)
}

type Connection interface {
	Id() string
	Send(packet.IPacket) error
	RemoteAddr() net.Addr
	Close() error
}

type Handler interface {
	OnConnect(Connection)
	OnDisConnect(Connection)
	OnReceived(Connection, packet.IPacket)
}

var (
	DefaultHandler = &handler{}
)

// handler default server handler
type handler struct {
}

func (h handler) OnConnect(conn Connection) {
}

func (h handler) OnDisConnect(conn Connection) {
}

func (h handler) OnReceived(conn Connection, packet packet.IPacket) {
}
