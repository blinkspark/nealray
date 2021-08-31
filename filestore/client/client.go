package main

import (
	"context"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

const (
	FILESTORE_PROTOCOL = "nealfree.cf/filestore/v0.1"
)

func main() {
	// use flag to get arguments
	// -d dial address (e.g. /ip4/127.0.0.1/tcp/4001)
	// -f filePath to read
	dialAddr := flag.String("d", "", "address to dial")
	filePath := flag.String("f", "", "file to read")
	flag.Parse()

	ctx := context.Background()
	h, err := libp2p.New(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	// connect to the peer
	if *dialAddr == "" {
		log.Panic("dial address is required")
	}
	addr, err := multiaddr.NewMultiaddr(*dialAddr)
	if err != nil {
		log.Panic(err)
	}
	// parse multiaddr to peer.AddrInfo
	pi, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Panic(err)
	}

	// dial to the peer
	if err := h.Connect(ctx, *pi); err != nil {
		log.Panic(err)
	}

	// read file
	if *filePath == "" {
		log.Panic("file path is required")
	}
	file, err := os.Open(*filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// send file to peer
	if s, err := h.NewStream(ctx, pi.ID, FILESTORE_PROTOCOL); err != nil {
		log.Panic(err)
	} else {
		defer s.Reset()
		_, err := s.Write([]byte{0x01})
		if err != nil {
			log.Panic(err)
		}
		// make a buffer to read file
		buf := make([]byte, 1024)
		fileNameLength := len(file.Name())
		binary.BigEndian.PutUint32(buf, uint32(fileNameLength))
		_, err = s.Write(buf[:4])
		if err != nil {
			log.Panic(err)
		}
		_, err = s.Write([]byte(file.Name()))
		if err != nil {
			log.Panic(err)
		}

		for {
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				log.Panic(err)
			}
			if n == 0 {
				break
			}
			if _, err := s.Write(buf[:n]); err != nil {
				log.Panic(err)
			}
		}
	}
}
