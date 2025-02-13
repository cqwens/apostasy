package main

import (
	"os"
	"os/exec"
	"text/template"
)

func installNomad() error {
	// First download the key to a file
	keyFile := "/tmp/hashicorp.gpg"
	downloadCmd := exec.Command("curl", "-fsSL", "-o", keyFile, "https://apt.releases.hashicorp.com/gpg")
	if err := downloadCmd.Run(); err != nil {
		return err
	}

	// Then add the key
	addKeyCmd := exec.Command("apt-key", "add", keyFile)
	if err := addKeyCmd.Run(); err != nil {
		return err
	}

	// Clean up
	if err := os.Remove(keyFile); err != nil {
		return err
	}

	cmds := []string{
		"apt-add-repository \"deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main\"",
		"apt-get update",
		"apt-get install -y nomad",
	}
	return runCommands(cmds)
}

func configureNomad(cfg Config) error {
	nomadConfig := `
datacenter = "dc1"
data_dir = "{{.NomadDataDir}}"
server {
  enabled = true
  bootstrap_expect = 1
}
client {
  enabled = true
}
ui {
  enabled = true
}`

	tmpl, err := template.New("nomad").Parse(nomadConfig)
	if err != nil {
		return err
	}

	f, err := os.Create("/etc/nomad.d/server.hcl")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, cfg); err != nil {
		return err
	}

	return runCommands([]string{
		"systemctl enable nomad",
		"systemctl start nomad",
	})
}
