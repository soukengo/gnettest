package main

import (
	log "github.com/golang/glog"
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/packet"
)

// DemoHandler default server DemoHandler
type DemoHandler struct {
}

func (h DemoHandler) OnConnect(conn network.Connection) {
	log.Infof("OnConnect id: %v", conn.Id())
}

func (h DemoHandler) OnDisConnect(conn network.Connection) {
	log.Infof("OnDisConnect id: %v", conn.Id())
}

func (h DemoHandler) OnReceived(conn network.Connection, p packet.IPacket) {
	log.Infof("OnReceived connId: %s: %v", conn.Id(), p)

}
