package main

import (
	"os"
	"text/template"
)

func installCaddy() error {
	cmds := []string{
		"apt install -y debian-keyring debian-archive-keyring apt-transport-https",
		"curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg",
		"curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list",
		"apt-get update",
		"apt-get install -y caddy",
	}
	return runCommands(cmds)
}

func configureCaddy(cfg Config) error {
	caddyConfig := `
{{.Domain}} {
    reverse_proxy localhost:4646
}`

	tmpl, err := template.New("caddy").Parse(caddyConfig)
	if err != nil {
		return err
	}

	f, err := os.Create("/etc/caddy/Caddyfile")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, cfg); err != nil {
		return err
	}

	return runCommands([]string{
		"systemctl enable caddy",
		"systemctl start caddy",
	})
}
