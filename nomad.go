package main

import (
	"os"
	"text/template"
)

func installNomad() error {
	cmds := []string{
		"curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -",
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
