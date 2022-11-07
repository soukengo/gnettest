package packet

import (
	"errors"
	"io"
)

var (
	ErrInvalidPacket = errors.New("invalid packet")
)

type IPacket interface {
	UnPackFrom(r io.Reader) (err error)
	PackTo(w io.Writer) (err error)
}

type IFactory interface {
	Offer(connId string) (p IPacket)
}
