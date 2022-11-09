package socket

import (
	"context"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/zhenjl/cityhash"
	"gnettest/api/comet"
	"gnettest/internal/comet/conf"
	"gnettest/internal/comet/model"
	"gnettest/pkg/server/socket"
	"gnettest/pkg/server/socket/network/tcp"
	"gnettest/pkg/server/socket/options"
	"gnettest/pkg/util/strings"
	"time"
)

type Server struct {
	c         *conf.Config
	buckets   []*Bucket
	bucketIdx uint32
	core      socket.Manager
	render    *Render
}

func NewServer(c *conf.Config) *Server {
	s := &Server{
		c:      c,
		render: &Render{},
	}
	// init bucket
	s.buckets = make([]*Bucket, c.Bucket.Size)
	s.bucketIdx = uint32(c.Bucket.Size)
	for i := 0; i < c.Bucket.Size; i++ {
		s.buckets[i] = NewBucket(c.Bucket)
	}
	mgr := socket.NewManager(s, options.WithPacketFactory(model.NewPacketFactory()))
	mgr.RegisterTCPServer(&tcp.Config{Address: c.Network.Tcp.Bind[0], Multicore: true, ReadBuf: c.Network.Tcp.ReadBuf, SendBuf: c.Network.Tcp.SendBuf})
	s.core = mgr
	return s
}

// Buckets return all buckets.
func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

// Bucket get the bucket by subkey.
func (s *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	return s.buckets[idx]
}

func (s *Server) Start(ctx context.Context) error {
	go s.autoPush()
	return s.core.Start()
}

func (s *Server) Stop(ctx context.Context) (err error) {
	return s.core.Close()
}

// 模拟向客户端推送数据
func (s *Server) autoPush() {
	size := 10
	list := &comet.MessageList{List: make([]*comet.Message, size)}
	for i := 0; i < size; i++ {
		list.List[i] = &comet.Message{Content: strings.RandStr(200)}
	}
	data, _ := proto.Marshal(list)
	log.Infof("data size: %d", len(data))
	p := model.NewPacket(comet.OpRawMessage, data)
	for {
		var pushReq = &model.PushGroupReq{
			GroupID: "10000",
			Packet:  p,
		}
		for _, bucket := range s.Buckets() {
			bucket.PushGroup(pushReq)
		}
		time.Sleep(time.Second)
	}
}
