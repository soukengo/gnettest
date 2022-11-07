package nbio

import (
	"bytes"
	log "github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/lesismal/nbio"
	"gnettest/pkg/server/socket/network"
	"gnettest/pkg/server/socket/packet"
	"net/url"
	"sync"
)

type nbioServer struct {
	id        string
	cfg       *Config
	endpoint  *url.URL
	parser    *packet.Parser
	handler   network.Handler
	eng       *nbio.Engine
	startOnce *sync.Once
}

func NewServer(cfg *Config, parser *packet.Parser) network.Server {
	g := nbio.NewEngine(nbio.Config{
		Network:            "tcp",
		Addrs:              []string{cfg.Address},
		ReadBufferSize:     cfg.ReadBuf,
		MaxWriteBufferSize: cfg.SendBuf,
	})
	s := &nbioServer{
		id:        uuid.New().String(),
		endpoint:  &url.URL{Scheme: "tcp", Path: cfg.Address},
		parser:    parser,
		handler:   network.DefaultHandler,
		cfg:       cfg,
		startOnce: new(sync.Once),
		eng:       g,
	}
	g.OnOpen(s.onOpen)
	g.OnClose(s.onClose)
	g.OnRead(s.onRead)
	//g.OnData(s.onData)
	return s

}

func (s *nbioServer) Id() string {
	return s.id
}
func (s *nbioServer) Start() (err error) {
	s.startOnce.Do(func() {
		go func() {
			err = s.eng.Start()
			if err != nil {
				panic(err)
			}
		}()
		log.Infof("started tcp server with endpoint %s", s.endpoint)
	})
	return
}

func (s *nbioServer) Close() error {
	s.eng.Stop()
	return nil
}

func (s *nbioServer) EndPoint() *url.URL {
	return s.endpoint
}
func (s *nbioServer) SetHandler(handler network.Handler) {
	s.handler = handler
}

func (s *nbioServer) onOpen(c *nbio.Conn) {
	conn := newConn(c, s.parser)
	c.SetSession(conn)
	s.handler.OnConnect(conn)
}

func (s *nbioServer) onClose(c *nbio.Conn, err error) {
	conn := s.parseConn(c)
	if conn == nil {
		return
	}
	conn.markClosed()
	s.handler.OnDisConnect(conn)
}

func (s *nbioServer) onRead(c *nbio.Conn) {
	conn := s.parseConn(c)
	if conn == nil {
		c.Close()
		return
	}

	for {
		p, err := conn.Read()
		if err != nil {
			if err != packet.ErrInvalidPacket {
				log.Errorf("reader error: %v", err)
			}
			return
		}
		s.handler.OnReceived(conn, p)
		//time.Sleep(time.Second)
	}
}

func (s *nbioServer) onData(c *nbio.Conn, data []byte) {
	conn := s.parseConn(c)
	if conn == nil {
		c.Close()
		return
	}
	p, err := s.parser.Parse(conn.id, bytes.NewBuffer(data))
	if err != nil {
		if err != packet.ErrInvalidPacket {
			log.Errorf("reader error: %v", err)
		}
		return
	}
	s.handler.OnReceived(conn, p)
}

func (s *nbioServer) parseConn(c *nbio.Conn) *nbioConn {
	conn, ok := c.Session().(*nbioConn)
	if !ok {
		return nil
	}
	return conn
}
