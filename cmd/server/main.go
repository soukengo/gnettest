package main

import (
	"flag"
	"gnettest/pkg/pprof"
	"gnettest/pkg/server/socket/network/gnet"
	"gnettest/pkg/server/socket/packet"
	"os"
	"os/signal"
	"syscall"
)

var (
	serverAddr = ":5311"
	pprofAddr  = ":5301"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "2")
}

func main() {
	flag.Parse()
	parser := packet.NewPacketParser(packet.NewSimplePacketFactory())
	s := gnet.NewServer(&gnet.Config{Address: serverAddr, SendBuf: 4096, ReadBuf: 4096, Multicore: true}, parser)
	//s := nbio.NewServer(&nbio.Config{Address: serverAddr, SendBuf: 4096, ReadBuf: 4096, Multicore: true}, parser)
	s.SetHandler(&DemoHandler{})
	s.Start()
	pprof.Start(pprofAddr)
	waitShutdown()
}

// waitShutdown 等待关闭服务信号
func waitShutdown() {
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
