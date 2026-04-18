package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/tuist/infra/internal/api"
)

func main() {
	var listenAddress string
	flag.StringVar(&listenAddress, "listen-address", ":8090", "Address for the public API server.")
	flag.Parse()

	// This server is the product-facing HTTP entrypoint for Infra.
	// Clients use it to interact with Infra without talking to the Kubernetes
	// API directly.
	server := api.New(api.Config{ListenAddress: listenAddress})
	log.Printf("starting infra API server on %s", listenAddress)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("infra API server exited with error: %v", err)
	}
}
