package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"

	"github.com/blinkspark/prototypes/blinkserver/config"
)

var (
	configPath *string
)

func init() {
	configPath = flag.String("c", "config.json", "-c config.json")
	flag.Parse()
}

func main() {
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Panic(err)
	}

	for _, s := range cfg.Servers {
		log.Printf("%#+v\n", s)
		server := http.Server{
			Addr: s.Addr,
		}
		if s.Root != "" {
			var fileServer http.Handler
			if s.DisableListing {
				fileServer = http.FileServer(&config.DisableListingFileSystem{FS: http.Dir(s.Root)})
			} else {
				fileServer = http.FileServer(http.Dir(s.Root))
			}
			server.Handler = fileServer
		}
		if s.To != "" {
			url, err := url.Parse(s.To)
			if err != nil {
				log.Panic(err)
			}
			server.Handler = httputil.NewSingleHostReverseProxy(url)
		}
		go func() {
			server.ListenAndServe()
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	log.Printf("%#+v\n", sig)
}
