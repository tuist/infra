package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/tuist/kubebox/internal/host"
)

func main() {
	var (
		listenAddress string
		mode          string
		backend       string
		hostID        string
	)

	flag.StringVar(&listenAddress, "listen-address", ":8091", "Address for the host HTTP server.")
	flag.StringVar(&mode, "mode", "local", "Agent mode: local, linux, or macos.")
	flag.StringVar(&backend, "backend", "local", "Backend runtime: local, cloud-hypervisor, or tart.")
	flag.StringVar(&hostID, "host-id", "local-dev-host", "Stable identifier for the host.")
	flag.Parse()

	a, err := host.New(host.Config{
		ListenAddress: listenAddress,
		Mode:          mode,
		Backend:       backend,
		HostID:        hostID,
	})
	if err != nil {
		log.Fatalf("unable to create host service: %v", err)
	}

	log.Printf("starting kubebox host service on %s (mode=%s backend=%s hostID=%s)", listenAddress, mode, backend, hostID)
	if err := a.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("host service exited with error: %v", err)
	}
}
