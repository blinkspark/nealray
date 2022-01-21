package main

import (
	"context"
	"flag"
	"log"

	"github.com/blinkspark/prototypes/filestore/node"
)

const FILE_PROTOCOL = "/filestore/0.1.0"

type Args struct {
	Port    int
	KeyPath string
}

var args Args

func init() {
	flag.StringVar(&args.KeyPath, "key", "priv.key", "path to private key")
	flag.StringVar(&args.KeyPath, "k", "priv.key", "path to private key")
	flag.IntVar(&args.Port, "p", 22233, "port to listen on")
	flag.IntVar(&args.Port, "port", 22233, "port to listen on")
	flag.Parse()
}

func main() {
	p2pNode, err := node.New(args.KeyPath, args.Port)
	if err != nil {
		log.Panic(err)
	}
	log.Println(p2pNode.Host.ID(), p2pNode.Host.Addrs())

	err = p2pNode.Bootstrap()
	if err != nil {
		log.Panic(err)
	}

	<-p2pNode.BootstrapDoneChan

	go func() {
		p2pNode.Discovery.Advertise(context.Background(), FILE_PROTOCOL)
		pc, err := p2pNode.Discovery.FindPeers(context.Background(), FILE_PROTOCOL)
		if err != nil {
			log.Panic(err)
		}
		for p := range pc {
			if p.ID == p2pNode.Host.ID() {
				continue
			}
			log.Println("Found peer:", p)
			p2pNode.Host.Connect(context.Background(), p)
		}
	}()

	// topic, err := p2pNode.Pubsub.Join(FILE_PROTOCOL)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// go func() {
	// 	for {
	// 		err := topic.Publish(context.Background(), []byte("hello world"))
	// 		if err != nil {
	// 			log.Println(err)
	// 			continue
	// 		}
	// 		time.Sleep(time.Second)
	// 	}
	// }()
	// go func() {
	// 	sub, err := topic.Subscribe()
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	for {
	// 		msg, err := sub.Next(context.Background())
	// 		if err != nil {
	// 			log.Println(err)
	// 			continue
	// 		}
	// 		log.Println("from:", msg.GetFrom())
	// 		log.Println("data:", string(msg.Data))
	// 	}
	// }()

	select {}
}

// func handleFileStream(s network.Stream) {
//
// }
