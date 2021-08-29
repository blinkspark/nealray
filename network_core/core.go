package networkcore

import (
	"context"
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/blinkspark/prototypes/util"
	"github.com/libp2p/go-libp2p"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/crypto"
)

type NetworkCore struct {
	core.Host
}

// NewNetworkCore creates a new NetworkCore
// keyPath is the path to the key file
// listenAddrs is a list of addresses to listen on e.g. /ip4/0.0.0.0/tcp/0 or /ip6/::/udp/0/quic
func NewNetworkCore(keyPath string, listenAddrs []string) (*NetworkCore, error) {
	// get key
	priv, err := getKey(keyPath)
	if err != nil {
		return nil, err
	}
	log.Printf("%#+v\n", priv)

	// add quic transport
	transport := libp2p.ChainOptions(
		libp2p.DefaultTransports,
		// libp2p.Transport(libp2pquic.NewTransport),
	)

	ctx := context.Background()
	h, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings(listenAddrs...),
		libp2p.Identity(priv),
		transport,
	)
	if err != nil {
		return nil, err
	}
	return &NetworkCore{h}, nil
}

// getKey gets the key from the key file
// or generates a new one if it doesn't exist
func getKey(keyPath string) (crypto.PrivKey, error) {
	var priv crypto.PrivKey
	var err error
	if util.PathExists(keyPath) {
		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
		priv, err = crypto.UnmarshalPrivateKey(keyData)
		if err != nil {
			return nil, err
		}
	} else {
		priv, _, err = crypto.GenerateEd25519Key(rand.Reader)
		keyData, err := crypto.MarshalPrivateKey(priv)
		if err != nil {
			return nil, err
		}

		dir, _ := path.Split(keyPath)
		if dir != "" {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}

		err = ioutil.WriteFile(keyPath, keyData, 0644)
		if err != nil {
			return nil, err
		}
	}
	return priv, err
}
