package main

import (
	"errors"
	"flag"
	"log"

	"github.com/blinkspark/prototypes/blinkcrypto/cryptor"
)

var (
	iFile   *string
	oFile   *string
	cfgFile *string
	d       *bool
	i       *bool
	t       *int
)

func init() {
	iFile = flag.String("i", "", "-i [input_file]")
	oFile = flag.String("o", "", "-o [output_file]")
	cfgFile = flag.String("c", "config.json", "-c [configfile.json]")
	d = flag.Bool("d", false, "-d #flag for decrypt")
	i = flag.Bool("init", false, "-init #flag for init a config file")
	t = flag.Int("t", 0, "-t 0 # 0 for chacha20poly1305, 1 for xchacha20poly1305")
	flag.Parse()
}

func main() {
	var (
		c   *cryptor.Cryptor
		err error
	)
	if *i {
		switch *t {
		case int(cryptor.CryptorChacha20poly1305):
			c, err = cryptor.NewChacha20Poly1305Cryptor()
		case int(cryptor.CryptorXChacha20poly1305):
			c, err = cryptor.NewXChacha20Poly1305Cryptor()
		default:
			err = errors.New("type error")
			log.Panic(err)
		}
		if err == nil {
			err = c.Save(*cfgFile)

			if err != nil {
				log.Panic(err)
			}
			log.Println("config file saved:", *cfgFile)
			return
		}
	} else {
		c, err = cryptor.LoadFromJson(*cfgFile)
	}
	if err != nil {
		log.Panic(err)
	}

	if *d {
		err = c.DecryptFile("main.go.enc", "main.go.dec")
		if err != nil {
			log.Panic(err)
		}
	} else {
		err = c.EncryptFile("main.go", "main.go.enc")
		if err != nil {
			log.Panic(err)
		}
	}
}
