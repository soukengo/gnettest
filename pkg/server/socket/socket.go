package socket

import (
	"gnettest/pkg/server/socket/network/tcp"
	"gnettest/pkg/server/socket/packet"
)

type Manager interface {
	Start() error
	Close() error
	Channel(connId string) Channel
	ServerRegister
}

type ServerRegister interface {
	RegisterTCPServer(cfg *tcp.Config)
}

type Listener interface {
	OnCreated(Channel)
	OnClosed(Channel)
	OnReceived(Channel, packet.IPacket)
}

type Channel interface {
	Id() string
	ClientIP() string
	Send(packet packet.IPacket) error
	Close() error
}

type receiver interface {
	Receive(p packet.IPacket)
}
