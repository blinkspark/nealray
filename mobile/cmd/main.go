package main

import (
	"flag"

	"github.com/blinkspark/prototypes/mobile/client"
	"github.com/blinkspark/prototypes/mobile/server"
)

var (
	DialAddr string
)

func init() {
	flag.StringVar(&DialAddr, "d", "", "-d /ip4/127.0.0.1/tcp/1234")
	flag.Parse()
}

func main() {
	if DialAddr != "" {
		client.StartClient(DialAddr)
	} else {
		server.StartServer()
	}
}
