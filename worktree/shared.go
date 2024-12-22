package worktree

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func runCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Command failed with exit code %d: %v\n", exitErr.ExitCode(), string(exitErr.Stderr))
		} else {
			log.Printf("Failed to run command: %v\n", err)
		}
	}
	return strings.TrimSpace(string(output))
}
