package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func configureSSH() error {
	return appendToFile("/etc/ssh/sshd_config", "PermitRootLogin yes")
}

func installDependencies() error {
	cmds := []string{
		"apt-get update",
		"apt-get install -y curl gnupg software-properties-common",
	}
	return runCommands(cmds)
}

func runCommands(cmds []string) error {
	for _, cmd := range cmds {
		args := strings.Fields(cmd)
		command := exec.Command(args[0], args[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return fmt.Errorf("command failed: %s: %v", cmd, err)
		}
	}
	return nil
}

func appendToFile(filename, text string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(text + "\n"); err != nil {
		return err
	}
	return nil
}
