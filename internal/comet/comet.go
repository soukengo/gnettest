package comet

import (
	"github.com/go-kratos/kratos/v2"
	"gnettest/internal/comet/conf"
	"gnettest/internal/comet/server/socket"
	"gnettest/pkg/pprof"
)

func New() (app *kratos.App) {
	cfg := conf.Default()
	app, err := wireApp(cfg)
	if err != nil {
		panic(err)
	}
	pprof.Start(":5555")
	return
}
func newApp(ss *socket.Server) *kratos.App {
	return kratos.New(
		kratos.Name("gnettest"),
		kratos.Version("1.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			ss,
		),
	)
}
