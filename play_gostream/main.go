package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	gostream "github.com/libp2p/go-libp2p-gostream"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	// parse command line options
	// -d dial to a peer addr (eg. /ip4/127.0.0.1/tcp/4001/p2p/QmAddr)
	// -l listen for incoming connections (eg. /ip4/127.0.0.1/tcp/4001)
	dialAddr := flag.String("d", "", "dial to peer address")
	listenAddr := flag.String("l", "", "listen for incoming connections")
	flag.Parse()

	if *dialAddr != "" {
		go runClient(*dialAddr)
	}

	if *listenAddr != "" {
		go runServer(*listenAddr)
	}
	select {}
}

// runClient runs a client for the given peer address.
func runClient(addr string) {
	// create a new libp2p node
	ctx := context.Background()
	client, err := libp2p.New(ctx)
	if err != nil {
		log.Panic(err)
	}

	// parse string to multiaddr
	parsedAddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Panic(err)
	}
	// parse multiaddr to peer.AddrInfo
	peerInfo, err := peer.AddrInfoFromP2pAddr(parsedAddr)
	if err != nil {
		log.Panic(err)
	}

	// dial to peer
	if err := client.Connect(ctx, *peerInfo); err != nil {
		log.Panic(err)
	}

	// create a new net.Conn
	conn, err := gostream.Dial(ctx, client, peerInfo.ID, "test")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	// write to the net.Conn
	conn.Write([]byte("hello world"))
}

// runServer runs a server for the given listen address.
func runServer(addr string) {
	// create a new libp2p node
	ctx := context.Background()
	laddr := libp2p.ListenAddrStrings(addr)
	server, err := libp2p.New(ctx, laddr)
	if err != nil {
		log.Panic(err)
	}
	log.Println("server listening on", server.Addrs())
	log.Println("server peer id", server.ID())

	listener, err := gostream.Listen(server, "test")
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			for {
				n, err := c.Read(buf)
				if err != nil {
					log.Println(err)
					return
				}
				log.Println(c.LocalAddr(), c.RemoteAddr(), string(buf[:n]))
			}
		}(conn)
	}
}
