package model

import (
	"bufio"
	"encoding/binary"
	"gnettest/pkg/server/socket/packet"
	"io"
)

const (
	packSize   = 4
	opSize     = 2
	codeSize   = 2
	headerSize = packSize + opSize + codeSize
	opOffset   = packSize
	codeOffset = opOffset + codeSize
)

type Packet struct {
	Op         uint16
	Code       uint16
	Body       []byte
	PacketData []byte
}

func (p *Packet) Encode() []byte {
	if len(p.PacketData) > 0 {
		return p.PacketData
	}
	p.PacketData = p.encode()
	return p.PacketData
}

func (p *Packet) Reset() {
	p.Op = 0
	p.Code = 0
	p.Body = nil
	p.PacketData = nil
}

func (p *Packet) UnPackFrom(r io.Reader) (err error) {
	bio, ok := r.(*bufio.Reader)
	if !ok {
		bio = bufio.NewReader(r)
	}
	packetLen := headerSize
	buf, _ := bio.Peek(packetLen)
	if len(buf) < packetLen {
		return packet.ErrInvalidPacket
	}
	packetLen = int(binary.BigEndian.Uint32(buf[:packSize]))
	p.Op = binary.BigEndian.Uint16(buf[opOffset:codeOffset])
	p.Code = binary.BigEndian.Uint16(buf[codeOffset:headerSize])

	buf, _ = bio.Peek(packetLen)
	if len(buf) < headerSize {
		return packet.ErrInvalidPacket
	}
	var body []byte
	body = buf[headerSize:]
	// discard
	_, _ = bio.Discard(packetLen)
	p.Body = body
	return
}

func (p *Packet) PackTo(w io.Writer) (err error) {
	packetData := p.PacketData
	if len(packetData) > 0 {
		_, err = w.Write(packetData)
		if err != nil {
			return
		}
		return
	}
	buf := p.encode()
	_, err = w.Write(buf)
	return
}

func (p *Packet) encode() []byte {
	bodySize := len(p.Body)
	packLength := headerSize + bodySize
	buf := make([]byte, packLength)
	binary.BigEndian.PutUint32(buf[0:], uint32(packLength))
	binary.BigEndian.PutUint16(buf[opOffset:], uint16(p.Op))
	binary.BigEndian.PutUint16(buf[codeOffset:], uint16(p.Code))
	copy(buf[headerSize:], p.Body)
	return buf
}

type PacketFactory struct {
	newFunc func() packet.IPacket
}

func (factory *PacketFactory) Offer(connId string) (p packet.IPacket) {
	return factory.newFunc()
}

func NewPacketFactory() packet.IFactory {
	return &PacketFactory{newFunc: func() packet.IPacket {
		return &Packet{}
	}}
}

func NewPacket(op uint16, body []byte) *Packet {
	return &Packet{Op: op, Body: body}
}
