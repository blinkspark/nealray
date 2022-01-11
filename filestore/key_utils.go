package main

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// getEd25519PrivateKey load key from file or generate new key and save it to file
func getEd25519PrivateKey(keyFile string) (crypto.PrivKey, error) {
	// Load key from file
	key, err := loadEd25519PrivateKey(keyFile)
	if err == nil {
		return key, nil
	}

	// Generate new key
	key, err = generateEd25519PrivateKey()
	if err != nil {
		return nil, err
	}

	// Save key to file
	err = saveEd25519PrivateKey(key, keyFile)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func loadEd25519PrivateKey(keyFile string) (crypto.PrivKey, error) {
	// Read key from file
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	// Decode key
	return crypto.UnmarshalPrivateKey(keyBytes)
}

func generateEd25519PrivateKey() (crypto.PrivKey, error) {
	// Generate new key
	key, _, err := crypto.GenerateEd25519Key(rand.Reader)
	return key, err
}

func saveEd25519PrivateKey(key crypto.PrivKey, keyFile string) error {
	// Encode key
	keyBytes, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return err
	}
	// Write key to file
	return os.WriteFile(keyFile, keyBytes, 0644)
}
