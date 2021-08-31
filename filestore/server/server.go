package main

import (
	"log"

	"github.com/blinkspark/prototypes/filestore"
)

// this is the server of file store

func main() {
	fstore, err := filestore.NewFSFileStore("data", "id.key", []string{"/ip4/0.0.0.0/tcp/33445"})
	if err != nil {
		log.Panic(err)
	}
	log.Println(fstore.ID())
	log.Println(fstore.Addrs())
	log.Printf("%#+v", fstore)
	select {}
}
