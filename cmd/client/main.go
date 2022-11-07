package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/golang/glog"
	"gnettest/pkg/server/socket/packet"
	"golang.org/x/net/context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	protocol   = "tcp"
	serverAddr = ":5311"
)

var (
	size = int64(1)
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
	for i := 0; i < int(size); i++ {
		newClient(ctx)
	}
	waitShutdown(func() {
		cancel()
	})
}

func newClient(ctx context.Context) {
	conn, err := net.Dial(protocol, serverAddr)
	if err != nil {
		log.Errorf("net.Dial err: %v", err)
		return
	}

	rr := bufio.NewReader(conn)
	// read
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				p := packet.NewSimplePacket(nil)
				err = p.UnPackFrom(rr)
				if err != nil {
					log.Errorf("UnPackFrom: err: %v", err)
					return
				}
				log.Infof("line: %v", p)
			}

		}

	}()

	for i := 0; i < 10; i++ {
		p := packet.NewSimplePacket([]byte(fmt.Sprintf("client body: %d", i)))
		err = p.PackTo(conn)
		if err != nil {
			log.Errorf("PackTo error: %v", err)
		}
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
