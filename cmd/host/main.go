package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/tuist/infra/internal/host"
)

func main() {
	var (
		listenAddress string
		mode          string
		vmRuntimes    string
		hostID        string
	)

	flag.StringVar(&listenAddress, "listen-address", ":8091", "Address for the host HTTP server.")
	flag.StringVar(&mode, "mode", "local", "Agent mode: local, linux, or macos.")
	flag.StringVar(&vmRuntimes, "vm-runtimes", "", "Comma-separated VM runtimes supported by this host, for example tart or cloud-hypervisor.")
	flag.StringVar(&hostID, "host-id", "local-dev-host", "Stable identifier for the host.")
	flag.Parse()

	a, err := host.New(host.Config{
		ListenAddress: listenAddress,
		Mode:          mode,
		VMRuntimes:    strings.Split(vmRuntimes, ","),
		HostID:        hostID,
	})
	if err != nil {
		log.Fatalf("unable to create host service: %v", err)
	}

	log.Printf("starting infra host service on %s (mode=%s vmRuntimes=%s hostID=%s)", listenAddress, mode, vmRuntimes, hostID)
	if err := a.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("host service exited with error: %v", err)
	}
}
