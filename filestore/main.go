package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func main() {
	key, err := getEd25519PrivateKey("priv.key")
	if err != nil {
		log.Panic(err)
	}

	node, err := libp2p.New(libp2p.Identity(key), libp2p.EnableAutoRelay())
	if err != nil {
		log.Panic(err)
	}
	log.Println(node.ID())

	dhtNode, err := dht.New(context.Background(), node)
	if err != nil {
		log.Panic(err)
	}
	log.Println(dhtNode)
	err = dhtNode.Bootstrap(context.Background())
	if err != nil {
		log.Panic(err)
	}

	pis, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		log.Panic(err)
	}

	for _, pi := range pis {
		log.Println(pi)
		err = node.Connect(context.Background(), pi)
		log.Println(err)
	}

	discNode := discovery.NewRoutingDiscovery(dhtNode)
	discNode.Advertise(context.Background(), "/ipfs/QmVr25JiAJiLVcHMLVu3q4FnQgCgtd1kxobUsr6qY5zGJY")

	discNode.FindPeers(context.Background(), "/ipfs/QmVr25JiAJiLVcHMLVu3q4FnQgCgtd1kxobUsr6qY5zGJY")
	select {}
}
