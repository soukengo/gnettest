package gnet

import (
	log "github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/panjf2000/gnet/v2"
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/packet"
	"sync"
	"time"
)

type gnetServer struct {
	id        string
	cfg       *Config
	parser    *packet.Parser
	handler   network.Handler
	eng       gnet.Engine
	startOnce *sync.Once
}

func NewServer(cfg *Config, parser *packet.Parser) *gnetServer {
	return &gnetServer{
		id:        uuid.New().String(),
		parser:    parser,
		handler:   network.DefaultHandler,
		cfg:       cfg,
		startOnce: new(sync.Once),
	}
}

func (s *gnetServer) Id() string {
	return s.id
}
func (s *gnetServer) Start() (err error) {
	s.startOnce.Do(func() {
		go func() {
			err = gnet.Run(s, s.cfg.Address,
				gnet.WithMulticore(s.cfg.Multicore),
				gnet.WithReuseAddr(true),
				//gnet.WithSocketSendBuffer(s.cfg.SendBuf),
				//gnet.WithSocketRecvBuffer(s.cfg.ReadBuf),
			)
			if err != nil {
				panic(err)
			}
		}()
		log.Infof("started tcp server with addr %s", s.cfg.Address)
	})
	return
}

func (s *gnetServer) Close() error {
	return nil
}

func (s *gnetServer) SetHandler(handler network.Handler) {
	s.handler = handler
}

func (s gnetServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.eng = eng
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 5)
	//		log.Infof("connection count: %d", s.eng.CountConnections())
	//	}
	//}()
	return
}

func (s *gnetServer) OnShutdown(eng gnet.Engine) {
	return
}

func (s *gnetServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	conn := newConn(c, s.parser)
	c.SetContext(conn)
	s.handler.OnConnect(conn)
	return
}

func (s *gnetServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	conn := s.parseConn(c)
	if conn == nil {
		return
	}
	conn.closed.CAS(false, true)
	s.handler.OnDisConnect(conn)
	return
}

func (s *gnetServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	conn := s.parseConn(c)
	if conn == nil {
		return gnet.Close
	}
	for {
		p, err := conn.Read()
		if err != nil {
			if err != packet.ErrInvalidPacket {
				log.Infof("reader error: %v", err)
			}
			return
		}
		s.handler.OnReceived(conn, p)
	}
}

func (s *gnetServer) OnTick() (delay time.Duration, action gnet.Action) {

	return
}

func (s *gnetServer) parseConn(c gnet.Conn) *tcpConn {
	conn, ok := c.Context().(*tcpConn)
	if !ok {
		return nil
	}
	return conn
}
