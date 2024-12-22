package worktree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// runCommand executes a command and returns its output

func CreateWorktreeCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create a new git worktree with a specified branch and base",
		ArgsUsage: "<branch-name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "base",
				Usage:   "Base branch to use (default: main)",
				Aliases: []string{"b"},
				Value:   "main",
			},
			&cli.BoolFlag{
				Name:    "upstream",
				Usage:   "Set upstream for the new branch",
				Aliases: []string{"u"},
				Value:   true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Check if branch name was provided
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("branch name is required")
			}

			branch := cmd.Args().First()
			base := cmd.String("base")
			upstream := cmd.Bool("upstream")

			// Get the current directory
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %v", err)
			}

			// Check if the current directory contains a .git file
			gitFile := filepath.Join(currentDir, ".git")
			if _, err := os.Stat(gitFile); os.IsNotExist(err) {
				return fmt.Errorf("no .git file found in the current directory: %v", err)
			}

			// Create the path for the new worktree
			worktreePath := filepath.Join(currentDir, branch)

			// Check if the directory already exists
			if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
				return fmt.Errorf("directory %s already exists", worktreePath)
			}

			fmt.Printf("Creating worktree at %s with branch %s based on %s\n", worktreePath, branch, base)

			// Create the worktree
			runCommand("git", "worktree", "add", worktreePath, "-b", branch, base)

			// Set upstream if the flag is true
			if upstream {
				// Change to the new worktree directory
				if err := os.Chdir(worktreePath); err != nil {
					return fmt.Errorf("failed to change to directory %s: %v", worktreePath, err)
				}

				fmt.Printf("Setting upstream branch and pushing to origin/%s\n", branch)
				runCommand("git", "push", "-u", "origin", branch)
			}

			fmt.Printf("Successfully created worktree for branch '%s'\n", branch)
			return nil
		},
	}
}
