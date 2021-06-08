package server

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	ps "github.com/libp2p/go-libp2p-pubsub"
)

func StartServer() {
	ctx := context.Background()
	h, err := libp2p.New(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	psNode, err := ps.NewGossipSub(ctx, h, ps.WithPeerExchange(true))
	if err != nil {
		log.Panic(err)
	}
	topic, err := psNode.Join("play")
	if err != nil {
		log.Panic(err)
	}
	defer topic.Close()

	sub, err := topic.Subscribe()
	if err != nil {
		log.Panic(err)
	}
	go func() {
		for {
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

	log.Println(h.Addrs())
	id := h.ID().Pretty()
	for _, addr := range h.Addrs() {
		log.Println(addr.String() + "/p2p/" + id)
	}

	log.Println(sub.Topic())
	select {}
}
