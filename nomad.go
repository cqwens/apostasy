package main

import (
	"os"
	"os/exec"
	"text/template"
)

func installNomad() error {
	// Download GPG key separately
	keyCmd := exec.Command("curl", "-fsSL", "https://apt.releases.hashicorp.com/gpg")
	addKeyCmd := exec.Command("apt-key", "add", "-")

	addKeyCmd.Stdin, _ = keyCmd.StdoutPipe()
	addKeyCmd.Stdout = os.Stdout
	addKeyCmd.Stderr = os.Stderr

	if err := addKeyCmd.Start(); err != nil {
		return err
	}
	if err := keyCmd.Run(); err != nil {
		return err
	}
	if err := addKeyCmd.Wait(); err != nil {
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
