package worktree

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func CloneCommand() *cli.Command {
	return &cli.Command{
		Name:  "clone",
		Usage: "Git clone bare repo with adjusted structure",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			repoURL := cmd.Args().First()
			basename := filepath.Base(repoURL)
			name := strings.TrimSuffix(basename, filepath.Ext(basename))
			createWorktrees := cmd.Bool("tree")

			if err := os.Mkdir(name, 0755); err != nil {
				log.Fatalf("Failed to create directory %s: %v", name, err)
			}
			if err := os.Chdir(name); err != nil {
				log.Fatalf("Failed to change to directory %s: %v", name, err)
			}

			runCommand("git", "clone", "--bare", repoURL, ".bare")

			// Create a .git file pointing to the .bare directory
			gitDir := "gitdir: ./.bare\n"
			if err := os.WriteFile(".git", []byte(gitDir), 0644); err != nil {
				log.Fatalf("Failed to write .git file: %v", err)
			}

			// Configure the remote fetch
			runCommand("git", "config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*")

			// Fetch all branches from the origin
			runCommand("git", "fetch", "origin")

			if createWorktrees {
				// Get remote branches and create worktrees
				output := runCommandOutput("git", "branch", "-r")
				branches := parseBranches(output)

				for _, branch := range branches {
					worktreeDir := strings.TrimPrefix(branch, "origin/")
					runCommand("git", "worktree", "add", worktreeDir, worktreeDir)
				}
			}

			return nil
		},

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "tree",
				Usage:   "create worktrees for all current branches",
				Aliases: []string{"t"},
				Value:   false,
			},
		},
	}
}

func runCommandOutput(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Command %s failed: %v", command, err)
	}
	return string(output)
}

func parseBranches(output string) []string {
	lines := strings.Split(output, "\n")
	var branches []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "origin/") {
			branches = append(branches, line)
		}
	}
	return branches
}
