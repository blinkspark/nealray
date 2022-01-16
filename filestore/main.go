package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	FILE_PROTOCOL = "/nealfree.cf/file/0.1.0"
	FILE_PUBSUB   = "/nealfree.cf/file/pubsub/0.1.0"
)

func main() {
	// key
	key, err := getEd25519KeyFromFile("priv.key")
	if err != nil {
		log.Panic(err)
	}

	// create libp2p node
	var dhtNode *dht.IpfsDHT
	node, err := libp2p.New(libp2p.Identity(key), libp2p.NATPortMap(), libp2p.EnableAutoRelay(), libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		dhtNode, err = dht.New(context.Background(), h)
		if err != nil {
			return nil, err
		}
		err = dhtNode.Bootstrap(context.Background())
		if err != nil {
			log.Panic(err)
		}
		return dhtNode, err
	}))
	if err != nil {
		log.Panic(err)
	}
	log.Println(node.ID())

	// bootstrap default dht servers
	pis, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		log.Panic(err)
	}
	for _, pi := range pis {
		log.Println(pi)
		err = node.Connect(context.Background(), pi)
		log.Println(err)
	}

	// gossip pubsub
	pubsubNode, err := pubsub.NewGossipSub(context.Background(), node)
	if err != nil {
		log.Panic(err)
	}

	topic, err := pubsubNode.Join(FILE_PUBSUB)
	if err != nil {
		log.Panic(err)
	}
	log.Println(topic)

	discNode := discovery.NewRoutingDiscovery(dhtNode)

	// advertise
	ttl, err := discNode.Advertise(context.Background(), FILE_PROTOCOL)
	log.Println(ttl, err)

	pchan, err := discNode.FindPeers(context.Background(), FILE_PROTOCOL)
	if err != nil {
		log.Panic(err)
	}
	for p := range pchan {
		log.Println(p)
	}

	select {}
}
