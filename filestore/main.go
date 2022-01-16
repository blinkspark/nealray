package main

import (
	"log"

	"github.com/blinkspark/prototypes/filestore/node"
)

func main() {
	p2pNode, err := node.New("priv.key")
	if err != nil {
		log.Panic(err)
	}
	log.Println(p2pNode.Host.ID(), p2pNode.Host.Addrs())
}
