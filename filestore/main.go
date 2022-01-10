package main

import "log"

func main() {
	priv, err := getEd25519KeyFromFile("priv.key")
	log.Println(priv, err)
}
