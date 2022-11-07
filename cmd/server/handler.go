package main

import (
	"fmt"
	log "github.com/golang/glog"
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/packet"
	"time"
)

// DemoHandler default server DemoHandler
type DemoHandler struct {
}

func (h DemoHandler) OnConnect(conn network.Connection) {
	log.Infof("OnConnect id: %v", conn.Id())
	//var ch *Channel
	//ch = newChannel(conn, 10, 10, func(p packet.IPacket) {
	//})
	//go pushMsg(ch)
}

func (h DemoHandler) OnDisConnect(conn network.Connection) {
}

func (h DemoHandler) OnReceived(conn network.Connection, packet packet.IPacket) {
	log.Infof("OnReceived connId: %s: %v", conn.Id(), packet)
}

// 模拟向客户端推送数据
func pushMsg(ch *Channel) {
	i := 0
	for {
		i++
		time.Sleep(time.Millisecond * 100)
		p := packet.NewSimplePacket([]byte(fmt.Sprintf("服务端主动推送的消息内容: %v", i)))
		ch.Send(p)
	}
}
