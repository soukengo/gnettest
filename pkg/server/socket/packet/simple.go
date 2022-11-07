package packet

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

type SimplePacket struct {
	bodyLen uint32
	body    []byte
}

func (p *SimplePacket) String() string {
	return fmt.Sprintf("bodyLen: %v,body: %v", p.bodyLen, string(p.body))
}

type simplePacketOptions struct {
	bodyLenSize int
	bodyOffset  int
}

var (
	simpleOpt = initSimplePacketOptions()
)

func initSimplePacketOptions() *simplePacketOptions {
	opt := &simplePacketOptions{
		bodyLenSize: 4,
	}
	opt.bodyOffset = opt.bodyLenSize
	return opt
}

func NewSimplePacket(body []byte) *SimplePacket {
	return &SimplePacket{
		bodyLen: uint32(len(body)),
		body:    body,
	}
}

func (p *SimplePacket) UnPackFrom(r io.Reader) (err error) {
	opt := simpleOpt
	bio, ok := r.(*bufio.Reader)
	if !ok {
		bio = bufio.NewReader(r)
	}
	packetLen := opt.bodyLenSize
	buf, _ := bio.Peek(packetLen)
	if len(buf) < packetLen {
		return ErrInvalidPacket
	}
	bodyLen := binary.BigEndian.Uint32(buf)
	// 整体包长度
	packetLen = packetLen + int(bodyLen)
	// 计算offset
	bodyOffset := opt.bodyOffset

	buf, _ = bio.Peek(packetLen)
	if len(buf) < packetLen {
		return ErrInvalidPacket
	}
	var body []byte
	if bodyLen > 0 {
		body = buf[bodyOffset:]
	}
	// discard
	_, _ = bio.Discard(packetLen)
	p.bodyLen = bodyLen
	p.body = body
	return
}

func (p *SimplePacket) PackTo(w io.Writer) (err error) {
	opt := simpleOpt
	buf := make([]byte, opt.bodyLenSize)
	binary.BigEndian.PutUint32(buf, p.bodyLen)
	_, err = w.Write(buf)
	if err != nil {
		return
	}
	if len(p.body) > 0 {
		_, err = w.Write(p.body)
		if err != nil {
			return
		}
	}
	return
}

type DefaultFactory struct {
	newFunc func() IPacket
}

func (factory *DefaultFactory) Offer(connId string) (p IPacket) {
	return factory.newFunc()
}

func NewSimplePacketFactory() IFactory {
	return &DefaultFactory{newFunc: func() IPacket {
		return &SimplePacket{}
	}}
}
