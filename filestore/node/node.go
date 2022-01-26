package node

import (
	"context"
	"log"
	"strconv"
	"strings"
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
	Host      host.Host
	DHT       *dht.IpfsDHT
	Pubsub    *pubsub.PubSub
	Discovery *discovery.RoutingDiscovery
	Protocol  string
	Logger    *log.Logger
}

func New(keyPath string, port int, protocol string, logger *log.Logger) (node *Node, err error) {
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
		logger.Println("Bootstrapped DHT")
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

	logger.Println("Created libp2p host")
	return &Node{
		Host:      h,
		DHT:       dhtNode,
		Pubsub:    pubsubNode,
		Discovery: disc,
		Protocol:  protocol,
		Logger:    logger,
	}, nil
}

func (n *Node) bootstrap() {
	pis, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		n.Logger.Panic(err)
	}

	for _, pi := range pis {
		go func(pi peer.AddrInfo) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			// n.Logger.Println("Bootstrapping DHT with", pi)
			n.Host.Connect(ctx, pi)
			// n.Logger.Println(err)
		}(pi)
	}
}

func (n *Node) Start() chan bool {
	done := make(chan bool, 1)
	go n.bootstrap()
	go n.advertise()
	go n.discover()
	return done
}

func (n *Node) advertise() {
	var err error
	var ttl time.Duration
	for {
		ttl, err = n.Discovery.Advertise(context.Background(), n.Protocol)
		if err != nil {
			n.Logger.Println(err)
			time.Sleep(time.Second)
			n.Logger.Println("Retrying to advertise")
			continue
		}
		n.Logger.Println("Advertising... ttl:", ttl)
		time.Sleep(ttl)
	}
}

func (n *Node) discover() {
	var peerChan <-chan peer.AddrInfo
	var err error
	for {
		peerChan, err = n.Discovery.FindPeers(context.Background(), n.Protocol)
		if err != nil {
			n.Logger.Println(err)
			time.Sleep(time.Second)
			n.Logger.Println("Retrying to find peers")
			continue
		}
		for pi := range peerChan {
			if pi.ID == n.Host.ID() {
				n.Logger.Println("Skipping self")
				continue
			}
			go func(pi peer.AddrInfo) {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				n.Logger.Println("Connecting to", pi)
				err = n.Host.Connect(ctx, pi)
				if err != nil {
					n.Logger.Println("Connection failed:", err)
				}
			}(pi)
		}
		time.Sleep(time.Second)
		n.Logger.Println("Finished discovering, restarting another round")
	}
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
