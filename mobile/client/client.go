package client

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
)

func StartClient(dialAddr string) {
	ctx := context.Background()
	h, err := libp2p.New(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	ps, err := pubsub.NewGossipSub(ctx, h, pubsub.WithPeerExchange(true))
	if err != nil {
		log.Panic(err)
	}

	topic, err := ps.Join("play")
	if err != nil {
		log.Panic(err)
	}

	defer topic.Close()
	// time.Sleep(1 * time.Second)

	go func() {
		for {
			sub, err := topic.Subscribe()
			if err != nil {
				log.Panic(err)
			}
			msg, err := sub.Next(ctx)
			if err != nil {
				log.Panic(err)
			}

			from, err := peer.IDFromBytes(msg.From)
			if err != nil {
				log.Panic(err)
			}
			log.Println(from.Pretty(), ":", string(msg.Data))
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			topic.Publish(ctx, []byte("Hello"))
		}
	}()

	ma, err := multiaddr.NewMultiaddr(dialAddr)
	if err != nil {
		log.Panic(err)
	}
	pi, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		log.Panic(err)
	}
	err = h.Connect(ctx, *pi)
	log.Println("connect :", err)

	log.Println(h.Addrs())
	id := h.ID().Pretty()
	for _, addr := range h.Addrs() {
		log.Println(addr.String() + "/p2p/" + id)
	}
	log.Println("Client")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// select {
	<-sigChan
	log.Println(h.Peerstore().Peers())
	// }
}
