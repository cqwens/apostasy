package main

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	NomadDataDir string
	Domain       string
}

func main() {
	if os.Geteuid() != 0 {
		log.Fatal("This program must be run as root")
	}

	if len(os.Args) != 2 {
		log.Fatal("Usage: apostasy example.com")
	}

	cfg := Config{
		NomadDataDir: "/opt/nomad",
		Domain:       os.Args[1],
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Configure SSH", configureSSH},
		{"Install Dependencies", installDependencies},
		{"Install Nomad", installNomad},
		{"Configure Nomad", func() error { return configureNomad(cfg) }},
		{"Install Caddy", installCaddy},
		{"Configure Caddy", func() error { return configureCaddy(cfg) }},
	}

	for _, step := range steps {
		fmt.Printf("Executing: %s\n", step.name)
		if err := step.fn(); err != nil {
			log.Fatalf("Error during %s: %v", step.name, err)
		}
	}
}
