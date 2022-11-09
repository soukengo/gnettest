package socket

import (
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/packet"
)

func (m *serverManager) OnConnect(conn network.Connection) {
	var ch Channel
	ch = newChannel(conn, m.opt.RecvQueueSize, m.opt.SendQueueSize, func(p packet.IPacket) {
		m.lis.OnReceived(ch, p)
	})
	m.Bucket(conn.Id()).PutChannel(ch)
	m.lis.OnCreated(ch)
}

func (m *serverManager) OnDisConnect(conn network.Connection) {
	channelId := conn.Id()
	ch := m.Channel(conn.Id())
	if ch == nil {
		return
	}
	m.Bucket(channelId).DelChannel(channelId)
	m.lis.OnClosed(ch)
}

func (m *serverManager) OnReceived(conn network.Connection, p packet.IPacket) {
	ch := m.Channel(conn.Id())
	if ch == nil {
		_ = conn.Close()
		return
	}
	if r, ok := ch.(receiver); ok {
		r.Receive(p)
	} else {
		m.lis.OnReceived(ch, p)
	}
}
