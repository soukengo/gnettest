package comet

const (
	OpHeartbeat      = uint16(2)
	OpHeartbeatReply = uint16(3)
	OpAuth           = uint16(7)
	OpAuthReply      = uint16(8)
	OpRawMessage     = uint16(1005)
)

const (
	ResultCodeSuccess = uint16(1000)
	ResultCodeFailed  = uint16(1001)
)
