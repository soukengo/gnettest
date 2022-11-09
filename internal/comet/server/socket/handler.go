package socket

import (
	"errors"
	"fmt"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"gnettest/api/comet"
	"gnettest/internal/comet/model"
	"gnettest/pkg/server/socket"
	"gnettest/pkg/server/socket/packet"
	"golang.org/x/net/context"
	"time"
)

func (s *Server) OnCreated(channel socket.Channel) {
	log.Infof("OnCreated:%v, %v", channel.Id(), channel.ClientIP())
	connId := channel.Id()
	// 超时未授权成功，则断开连接
	time.AfterFunc(s.c.Network.HandshakeTimeout, func() {
		var ch = s.Bucket(connId).Channel(connId)
		if ch == nil {
			_ = channel.Close()
			return
		}
	})

}

func (s *Server) OnClosed(channel socket.Channel) {
	log.Infof("OnClosed:%v, %v", channel.Id(), channel.ClientIP())
	connId := channel.Id()
	var ch = s.Bucket(connId).Channel(connId)
	if ch == nil {
		return
	}
	_ = s.Disconnect(context.TODO(), ch)
}

func (s *Server) OnReceived(channel socket.Channel, packet packet.IPacket) {
	//log.Infof("OnReceived: %v", packet)
	connId := channel.Id()
	ctx := context.TODO()
	req, ok := packet.(*model.Packet)
	if !ok {
		log.Infof("packet type is not supported")
		channel.Close()
		return
	}
	resp := &model.Packet{}
	var ch = s.Bucket(connId).Channel(connId)
	if req.Op != comet.OpAuth && ch == nil {
		s.render.OutputError(channel, resp, errors.New("未授权访问"))
		_ = channel.Send(resp)
		channel.Close()
		return
	}
	var (
		data proto.Message
		err  error
	)
	switch req.Op {
	case comet.OpAuth:
		if ch == nil {
			ch = NewChannel(channel)
		}
		data, err = s.Auth(ctx, ch, req)
		resp.Op = comet.OpAuthReply

	case comet.OpHeartbeat:
		data, err = s.Heartbeat(ctx, ch, req)
		resp.Op = comet.OpHeartbeatReply
	default:
		err = fmt.Errorf("不支持的操作类型: %v", req.Op)
	}
	if err != nil {
		s.render.OutputError(channel, resp, err)
		return
	}
	s.render.OutputSuccess(channel, resp, data)
}
