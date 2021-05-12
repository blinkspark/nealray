package cryptor

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/chacha20poly1305"
)

type CryptorType int

const (
	CryptorChacha20poly1305 CryptorType = iota
	CryptorXChacha20poly1305
)

const DefaultBufferSize = 1 << 20

type Cryptor struct {
	cipher.AEAD `json:"-"`
	//types support now chacha20poly1305,xchacha20poly1305
	Type CryptorType `josn:"type"`
	Key  []byte      `josn:"key"`
}

func (c *Cryptor) Save(fname string) error {
	d, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(fname, d, 0666)
}

func (c *Cryptor) EncryptFile(fname string, outFName string) error {
	buffer := make([]byte, DefaultBufferSize)

	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	of, err := os.Create(outFName)
	if err != nil {
		return err
	}

	nonce := make([]byte, c.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	_, err = of.Write(nonce)
	if err != nil {
		return err
	}

	for {
		n, err := f.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		out := c.Seal(nil, nonce, buffer[:n], nil)
		_, err = of.Write(out)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cryptor) DecryptFile(fname string, outFName string) error {
	buffer := make([]byte, DefaultBufferSize)

	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	of, err := os.Create(outFName)
	if err != nil {
		return err
	}

	nonce := make([]byte, c.NonceSize())
	_, err = f.Read(nonce)
	if err != nil {
		return err
	}
	for {
		n, err := f.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		out, err := c.Open(nil, nonce, buffer[:n], nil)
		if err != nil {
			return err
		}

		_, err = of.Write(out)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewChacha20Poly1305Cryptor() (c *Cryptor, err error) {
	key := make([]byte, chacha20poly1305.KeySize)
	_, err = rand.Read(key)
	if err != nil {
		return nil, err
	}
	ciph, err := chacha20poly1305.New(key)
	c = &Cryptor{
		AEAD: ciph,
		Key:  key,
		Type: CryptorChacha20poly1305,
	}
	return
}

func NewXChacha20Poly1305Cryptor() (c *Cryptor, err error) {
	key := make([]byte, chacha20poly1305.KeySize)
	_, err = rand.Read(key)
	if err != nil {
		return nil, err
	}
	ciph, err := chacha20poly1305.NewX(key)
	c = &Cryptor{
		AEAD: ciph,
		Key:  key,
		Type: CryptorXChacha20poly1305,
	}
	return
}

func LoadFromJson(fname string) (c *Cryptor, err error) {
	c = &Cryptor{}
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	switch c.Type {
	case CryptorChacha20poly1305:
		c.AEAD, err = chacha20poly1305.New(c.Key)
		return
	case CryptorXChacha20poly1305:
		c.AEAD, err = chacha20poly1305.NewX(c.Key)
		return
	default:
		return nil, errors.New("type error")
	}
}
