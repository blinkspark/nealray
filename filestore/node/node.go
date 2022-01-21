package node

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Node struct {
	Host              host.Host
	DHT               *dht.IpfsDHT
	Pubsub            *pubsub.PubSub
	Discovery         *discovery.RoutingDiscovery
	BootstrapDoneChan chan bool
}

func New(keyPath string, port int) (node *Node, err error) {
	var h host.Host
	var dhtNode *dht.IpfsDHT
	var pubsubNode *pubsub.PubSub

	priv, err := ed25519KeyFromFile(keyPath)
	if err != nil {
		return nil, err
	}

	h, err = libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrStrings(listenAddrs(port)...), libp2p.EnableAutoRelay(), libp2p.EnableHolePunching(), libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
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

	disc := discovery.NewRoutingDiscovery(dhtNode)

	log.Println("Created libp2p host")
	return &Node{
		Host:              h,
		DHT:               dhtNode,
		Pubsub:            pubsubNode,
		Discovery:         disc,
		BootstrapDoneChan: make(chan bool),
	}, nil
}

func (n *Node) Bootstrap() error {
	pis, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, pi := range pis {
		wg.Add(1)
		go func(pi peer.AddrInfo) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			// log.Println("Bootstrapping DHT with", pi)
			err = n.Host.Connect(ctx, pi)
			// log.Println(err)
		}(pi)
	}
	go func() {
		wg.Wait()
		n.BootstrapDoneChan <- true
	}()
	return nil
}

func listenAddrs(port int) []string {
	portString := strconv.Itoa(port)
	template := []string{
		"/ip4/0.0.0.0/tcp/{PORT}",
		"/ip4/0.0.0.0/tcp/{PORT}/ws",
		"/ip4/0.0.0.0/udp/{PORT}/quic",
		"/ip6/::/tcp/{PORT}",
		"/ip6/::/tcp/{PORT}/ws",
		"/ip6/::/udp/{PORT}/quic",
	}
	ret := []string{}
	for _, t := range template {
		ret = append(ret, strings.ReplaceAll(t, "{PORT}", portString))
	}
	return ret
}
