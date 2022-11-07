package packet

import (
	"io"
)

type Parser struct {
	factory IFactory
}

func NewPacketParser(factory IFactory) *Parser {
	return &Parser{factory: factory}
}

func (parser *Parser) Parse(connId string, reader io.Reader) (p IPacket, err error) {
	packet := parser.factory.Offer(connId)
	err = packet.UnPackFrom(reader)
	if err != nil {
		return
	}
	p = packet
	return
}
