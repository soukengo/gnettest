package server

import (
	"github.com/google/wire"
	"gnettest/internal/comet/server/socket"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(socket.NewServer)
