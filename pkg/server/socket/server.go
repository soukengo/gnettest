package socket

import (
	"gnettest/pkg/server/socket/network/tcp"
	"gnettest/pkg/server/socket/network/tcp/gnet"
)

func (m *serverManager) RegisterTCPServer(cfg *tcp.Config) {
	//s := nbio.NewServer(cfg, m.opt.Parser)
	s := gnet.NewServer(cfg, m.opt.Parser)
	m.register(s)
	return
}
