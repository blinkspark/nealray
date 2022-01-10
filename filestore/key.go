package main

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p-core/crypto"
)

func generateEd25519Key() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	return priv, err
}

func saveKeyToFile(key crypto.PrivKey, path string) error {
	data, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadKeyFromFile(path string) (crypto.PrivKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(data)
}

// getEd25519KeyFromFile get key from file or generate new one and save it to file
func getEd25519KeyFromFile(path string) (crypto.PrivKey, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		key, err := generateEd25519Key()
		if err != nil {
			return nil, err
		}
		err = saveKeyToFile(key, path)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return loadKeyFromFile(path)
}
