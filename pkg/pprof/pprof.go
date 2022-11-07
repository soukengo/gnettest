package pprof

import (
	log "github.com/golang/glog"
	"net/http"
	"net/http/pprof"
)

func register(handler *http.ServeMux) {
	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func Start(addr string) {
	mux := http.NewServeMux()
	register(mux)
	go func() {
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			panic(err)
		}
	}()
	log.Infof("pprof listening on %s", addr)
}
