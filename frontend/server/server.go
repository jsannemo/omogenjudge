package main

import (
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/google/logger"
)

var (
	webAddress = flag.String("frontend_listen", "127.0.0.1:61814", "The listen address for the frontend server")
)

func main() {
	logger.Init("omogenjudge-frontend", true, false, ioutil.Discard)
	logger.Fatal(http.ListenAndServe(*webAddress, configureRouter()))
}
