package node

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Node struct {
	Host   host.Host
	DHT    *dht.IpfsDHT
	Pubsub *pubsub.PubSub
}

func New(keyPath string) (node *Node, err error) {
	var h host.Host
	var dhtNode *dht.IpfsDHT
	var pubsubNode *pubsub.PubSub

	priv, err := ed25519KeyFromFile(keyPath)
	if err != nil {
		return nil, err
	}

	h, err = libp2p.New(libp2p.Identity(priv), libp2p.EnableAutoRelay(), libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		dhtNode, err = dht.New(context.Background(), h)
		if err != nil {
			return nil, err
		}
		err = dhtNode.Bootstrap(context.Background())
		if err != nil {
			return nil, err
		}
		log.Println("Bootstrapped DHT")
		return dhtNode, nil
	}))
	if err != nil {
		return nil, err
	}

	pubsubNode, err = pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		return nil, err
	}

	log.Println("Created libp2p host")
	return &Node{
		Host:   h,
		DHT:    dhtNode,
		Pubsub: pubsubNode,
	}, nil
}

func (n *Node) Bootstrap() error {
	pis, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		return err
	}

	for _, pi := range pis {
		go func(pi peer.AddrInfo) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			log.Println("Bootstrapping DHT with", pi)
			err = n.Host.Connect(ctx, pi)
			log.Println(err)
		}(pi)
	}
	return nil
}
