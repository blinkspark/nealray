package main

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	"github.com/blinkspark/prototypes/filestore/node"
)

const FILE_PROTOCOL = "/filestore/0.1.0"

type Args struct {
	Port    int
	KeyPath string
	Verbose bool
	LogPath string
}

var args Args

func init() {
	flag.StringVar(&args.KeyPath, "key", "priv.key", "path to private key")
	flag.StringVar(&args.KeyPath, "k", "priv.key", "path to private key")
	flag.IntVar(&args.Port, "p", 22233, "port to listen on")
	flag.IntVar(&args.Port, "port", 22233, "port to listen on")
	flag.BoolVar(&args.Verbose, "v", false, "verbose")
	flag.BoolVar(&args.Verbose, "verbose", false, "verbose")
	flag.StringVar(&args.LogPath, "logpath", "", "path to log file")

	flag.Parse()
}

func main() {
	timeStr := time.Now().Format("20060102150405")
	// log.New()
	logger := log.Default()
	if !args.Verbose {
		_, err := os.Stat(args.LogPath)
		if err == nil {
			err = os.MkdirAll(args.LogPath, 0755)
			if err != nil {
				log.Panic(err)
			}
		}
		f, err := os.Create(path.Join(args.LogPath, timeStr+".log"))
		if err != nil {
			log.Panic(err)
		}
		logger = log.New(f, "", log.LstdFlags)
	}

	p2pNode, err := node.New(args.KeyPath, args.Port, FILE_PROTOCOL, logger)
	if err != nil {
		logger.Panic(err)
	}
	logger.Println(p2pNode.Host.ID(), p2pNode.Host.Addrs())

	<-p2pNode.Start()
}
