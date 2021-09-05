package main

// filestore use libp2p as transport

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"

	"github.com/blinkspark/prototypes/util/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
)

type Config struct {
	PrivateKey  string
	ListenAddrs []string
}

func main() {
	// parse flags
	const defaultConfigPath = "config.json"
	var defaultListenAddrs = []string{"/ip4/0.0.0.0/tcp/12233"}
	var (
		newConfig  = flag.Bool("new_cfg", false, "create new config")
		configPath = flag.String("c", defaultConfigPath, "config file path")
		cfg        Config
	)
	flag.Parse()

	if *newConfig {
		// generate private key
		priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			log.Panic(err)
		}
		privData, err := crypto.MarshalPrivateKey(priv)
		if err != nil {
			log.Panic(err)
		}
		cfg.PrivateKey = base64.StdEncoding.EncodeToString(privData)
		cfg.ListenAddrs = defaultListenAddrs
		// save config
		if err := config.SaveConfig(defaultConfigPath, cfg); err != nil {
			log.Panic(err)
		}
		log.Println("new config generated")
	} else {
		// read config
		err := config.ReadConfig(*configPath, &cfg)
		if err != nil {
			log.Panic(err)
		}
	}

	// new libp2p host
	ctx := context.Background()
	privData, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	if err != nil {
		log.Panic(err)
	}
	priv, err := crypto.UnmarshalPrivateKey(privData)
	if err != nil {
		log.Panic(err)
	}
	h, err := libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(cfg.ListenAddrs...))
	if err != nil {
		log.Panic(err)
	}

	log.Println(h.ID())
	log.Println(h.Addrs())
}
