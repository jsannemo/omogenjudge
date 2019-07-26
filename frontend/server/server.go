package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	webAddress = flag.String("frontend_listen", "127.0.0.1:61813", "The listen address for the frontend server")
)

func main() {
	log.Fatal(http.ListenAndServe(*webAddress, configureRouter()))
}
