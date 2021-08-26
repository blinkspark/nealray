package filestore

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

type FileStore interface {
	core.Host
}

type FSFileStore struct {
	core.Host
	StorePath string
	KeyPath   string
}

func NewFSFileStore(keyPath, storePath string, listenAddrs []string) (*FSFileStore, error) {
	log.Println("hello")
	var fstore *FSFileStore = &FSFileStore{
		StorePath: storePath,
		KeyPath:   keyPath,
	}
	ctx := context.Background()

	priv, err := getKey(keyPath)
	if err != nil {
		return nil, err
	}

	fstore.Host, err = libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(listenAddrs...),
	)
	if err != nil {
		return nil, err
	}

	if !util.PathExists(storePath) {
		err = os.Mkdir(storePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return fstore, nil
}

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
