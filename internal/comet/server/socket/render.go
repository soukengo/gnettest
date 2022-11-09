package socket

import (
	"github.com/golang/protobuf/proto"
	"gnettest/api/comet"
	"gnettest/internal/comet/model"
	"gnettest/pkg/server/socket"
)

type Render struct {
}

func (r *Render) OutputSuccess(ch socket.Channel, proto *model.Packet, data proto.Message) {
	r.output(ch, proto, comet.ResultCodeSuccess, data)
}

func (r *Render) OutputError(ch socket.Channel, proto *model.Packet, err error) {
	if err == nil {
		r.output(ch, proto, comet.ResultCodeSuccess, &comet.Empty{})
		return
	}
	r.output(ch, proto, comet.ResultCodeFailed, &comet.Error{Code: uint32(comet.ResultCodeFailed), Message: err.Error()})
}

func (r *Render) output(ch socket.Channel, p *model.Packet, code uint16, data proto.Message) {
	p.Code = code
	p.Body, _ = r.encodeBody(data)
	_ = ch.Send(p)
}

func (r *Render) encodeBody(data proto.Message) ([]byte, error) {
	return proto.Marshal(data)
}
