package filestore

import (
	core "github.com/libp2p/go-libp2p-core"
)

type FSFileStore struct {
	core.Host
	StorePath string
	KeyPath   string
}
