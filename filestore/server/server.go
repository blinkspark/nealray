package main

import (
	"log"
	"net/http"

	"github.com/blinkspark/prototypes/filestore"
)

func main() {
	fstore, err := filestore.NewFileStore("data")
	if err != nil {
		log.Panic(err)
	}

	r := fstore.Router()
	http.ListenAndServe(":8888", r)
}
