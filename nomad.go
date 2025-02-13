package main

import (
	"os"
	"os/exec"
	"text/template"
)

func installNomad() error {
	// Install prerequisites
	cmds := []string{
		"apt-get update",
		"apt-get install -y wget gpg coreutils",
	}
	if err := runCommands(cmds); err != nil {
		return err
	}

	// Download and setup GPG key
	keyCmd := exec.Command("wget", "-O-", "https://apt.releases.hashicorp.com/gpg")
	keyFile := "/usr/share/keyrings/hashicorp-archive-keyring.gpg"
	gpgCmd := exec.Command("gpg", "--dearmor", "-o", keyFile)

	gpgCmd.Stdin, _ = keyCmd.StdoutPipe()
	if err := gpgCmd.Start(); err != nil {
		return err
	}
	if err := keyCmd.Run(); err != nil {
		return err
	}
	if err := gpgCmd.Wait(); err != nil {
		return err
	}

	// Add repository
	repoLine := "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
	repoFile := "/etc/apt/sources.list.d/hashicorp.list"
	if err := os.WriteFile(repoFile, []byte(repoLine+"\n"), 0644); err != nil {
		return err
	}

	// Install Nomad
	return runCommands([]string{
		"apt-get update",
		"apt-get install -y nomad",
	})
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
