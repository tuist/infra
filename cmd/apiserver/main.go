package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/tuist/kubebox/internal/apiserver"
)

func main() {
	var listenAddress string
	flag.StringVar(&listenAddress, "listen-address", ":8090", "Address for the public API server.")
	flag.Parse()

	server := apiserver.New(apiserver.Config{ListenAddress: listenAddress})
	log.Printf("starting kubebox API server on %s", listenAddress)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("kubebox API server exited with error: %v", err)
	}
}

