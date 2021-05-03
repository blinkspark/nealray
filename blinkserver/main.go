package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"

	"github.com/NYTimes/gziphandler"
	"github.com/blinkspark/prototypes/blinkserver/config"
	"golang.org/x/net/webdav"
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
		// file server
		if s.Root != "" {
			var fileServer http.Handler
			if s.DisableListing {
				fileServer = http.FileServer(&config.DisableListingFileSystem{FS: http.Dir(s.Root)})
			} else {
				fileServer = http.FileServer(http.Dir(s.Root))
			}
			if s.EnableGzip {
				fileServer = gziphandler.GzipHandler(fileServer)
			}
			server.Handler = fileServer
		}
		// reverse proxy
		if s.To != "" {
			url, err := url.Parse(s.To)
			if err != nil {
				log.Panic(err)
			}
			var handler http.Handler
			handler = httputil.NewSingleHostReverseProxy(url)
			if s.EnableGzip {
				handler = gziphandler.GzipHandler(handler)
			}
			server.Handler = handler
		}
		// webdav
		if len(s.WebDAVs) > 0 {
			davmux := http.ServeMux{}
			for _, dav := range s.WebDAVs {
				h := webdav.Handler{
					Prefix:     dav.Prefix,
					FileSystem: webdav.Dir(dav.Root),
					LockSystem: webdav.NewMemLS(),
				}
				davmux.HandleFunc(dav.Prefix, func(w http.ResponseWriter, r *http.Request) {
					if len(dav.Users) > 0 {
						uname, passwd, _ := r.BasicAuth()
						for _, user := range dav.Users {
							if user.User == uname && user.Password == passwd {
								w.Header().Set("Timeout", "99999999")
								h.ServeHTTP(w, r)
								return
							}
						}
						w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`)
						w.WriteHeader(401)
						w.Write([]byte("401 Unauthorized\n"))
					} else {
						w.Header().Set("Timeout", "99999999")
						h.ServeHTTP(w, r)
					}
				})
			}
			var handler http.Handler
			handler = &davmux
			if s.EnableGzip {
				handler = gziphandler.GzipHandler(handler)
			}
			server.Handler = handler
		}
		go func(s config.ServerConfig) {
			if s.Cert != "" && s.Key != "" {
				server.ListenAndServeTLS(s.Cert, s.Key)
			} else {
				server.ListenAndServe()
			}
		}(s)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	log.Printf("server stoped, received signal:%#+v\n", sig)
}
