package main

import (
	"bufio"
	"flag"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"gnettest/api/comet"
	"gnettest/internal/comet/model"
	"golang.org/x/net/context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	protocol   = "tcp"
	serverAddr = ":5311"
)

var (
	size = int64(100)
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "2")
}

func main() {
	flag.Parse()
	if len(os.Args) > 1 {
		size, _ = strconv.ParseInt(os.Args[1], 10, 64)
	}
	ctx, cancel := context.WithCancel(context.TODO())
	for i := int64(0); i < size; i++ {
		uid := strconv.FormatInt(i, 10)
		newClient(ctx, uid)
	}
	waitShutdown(func() {
		cancel()
	})
}

func newClient(ctx context.Context, uid string) {
	conn, err := net.Dial(protocol, serverAddr)
	if err != nil {
		log.Errorf("net.Dial err: %v", err)
		return
	}
	rr := bufio.NewReader(conn)
	authed := make(chan struct{})
	// read
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				p := model.NewPacket(0, nil)
				err = p.UnPackFrom(rr)
				if err != nil {
					log.Errorf("UnPackFrom: err: %v", err)
					return
				}
				//log.Infof("line: %v", p)
				if p.Op == comet.OpAuthReply {
					authed <- struct{}{}
				}
			}

		}

	}()
	// heartbeat
	go func() {

		select {
		case <-ctx.Done():
			return
		case <-authed:
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				hbReq := &comet.HeartBeatRequest{}
				data, _ := proto.Marshal(hbReq)
				p := model.NewPacket(comet.OpHeartbeat, data)
				err := p.PackTo(conn)
				if err != nil {
					log.Errorf("PackTo error: %v", err)
					return
				}
			}
			time.Sleep(time.Second * 10)
		}
	}()
	authReq := &comet.AuthRequest{
		Uid: uid,
	}
	data, _ := proto.Marshal(authReq)
	p := model.NewPacket(comet.OpAuth, data)
	err = p.PackTo(conn)
	if err != nil {
		log.Errorf("PackTo error: %v", err)
	}

}

// waitShutdown 等待关闭服务信号
func waitShutdown(callback func()) {
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			callback()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
