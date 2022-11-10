package socket

import (
	"context"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"gnettest/api/comet"
	"gnettest/internal/comet/model"
	"time"
)

func (s *Server) Auth(ctx context.Context, ch *Channel, p *model.Packet) (res *comet.AuthResponse, err error) {
	log.Infof("Auth Start Ip: %v, Key: %v", ch.ClientIp(), ch.Key)
	req := comet.AuthRequest{}
	err = s.bindRequest(p.Body, &req)
	if err != nil {
		return
	}
	// 模拟rpc请求
	time.Sleep(time.Millisecond * 50)
	ch.Uid = req.Uid

	ch.Enable(time.Minute*5, func() {
		_ = s.Disconnect(context.TODO(), ch)
	})
	err = s.Bucket(ch.Key).Put(ch)
	// 加入组，用于接收广播推送消息
	ch.JoinGroup("10000")
	log.Infof("Auth Success Ip: %v, Key: %v", ch.ClientIp(), ch.Key)
	return &comet.AuthResponse{}, nil
}
func (s *Server) Heartbeat(ctx context.Context, ch *Channel, p *model.Packet) (res *comet.HeartBeatResponse, err error) {
	ch.Renew()
	res = &comet.HeartBeatResponse{}
	return
}

func (s *Server) Disconnect(ctx context.Context, ch *Channel) (err error) {
	s.Bucket(ch.Key).Del(ch)
	_ = ch.Close()
	return
}

func (s *Server) bindRequest(data []byte, req proto.Message) error {
	return proto.Unmarshal(data, req)
}
