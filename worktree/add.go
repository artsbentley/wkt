package worktree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

// AddWorktreeCommand returns a CLI command to create a new worktree as a sibling to the current worktree.
func AddWorktreeCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new git worktree with a specified branch and base",
		ArgsUsage: "<branch-name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "base",
				Usage:   "Base branch to use",
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
			// Validate branch name argument.
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("branch name is required")
			}

			branch := cmd.Args().First()
			base := cmd.String("base")
			upstream := cmd.Bool("upstream")

			// Get the top-level directory of the current worktree.
			repoRoot := runCommand("git", "rev-parse", "--show-toplevel")
			if repoRoot == "" {
				return fmt.Errorf("failed to get repository top-level directory")
			}

			// Use the parent of the current worktree as the base location.
			// (e.g. if your current worktree is in "main", the new one will be a sibling directory.)
			parentDir := filepath.Dir(repoRoot)

			// If the branch name contains slashes, sanitize it by replacing them with dashes.
			// This ensures the worktree is created as a single directory.
			sanitizedBranch := strings.ReplaceAll(branch, "/", "-")
			worktreePath := filepath.Join(parentDir, sanitizedBranch)

			// Check if the target directory already exists.
			if _, err := os.Stat(worktreePath); err == nil {
				return fmt.Errorf("directory %s already exists", worktreePath)
			}

			fmt.Printf("Creating worktree at %s with branch %s based on %s\n", worktreePath, branch, base)

			// Create the worktree. (The branch name passed to git remains unchanged.)
			runCommand("git", "worktree", "add", worktreePath, "-b", branch, base)

			// Set upstream if requested.
			if upstream {
				// Change to the new worktree directory.
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
