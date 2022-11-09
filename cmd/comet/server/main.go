package main

import (
	"flag"
	"gnettest/internal/comet"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "2")
}

func main() {
	flag.Parse()
	srv := comet.New()
	srv.Run()
}
