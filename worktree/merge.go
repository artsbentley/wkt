package worktree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func MergeWorktreeCommand() *cli.Command {
	return &cli.Command{
		Name:  "merge",
		Usage: "Merge or rebase the current worktree with the main branch",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "target",
				Usage:   "Target branch to merge with",
				Aliases: []string{"t"},
				Value:   "main",
			},
			&cli.BoolFlag{
				Name:    "rebase",
				Usage:   "Use rebase instead of merge",
				Aliases: []string{"r"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "push",
				Usage:   "Push changes after successful merge/rebase",
				Aliases: []string{"p"},
				Value:   true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Get the current directory
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %v", err)
			}

			// Verify we're in a git repository
			if _, err := os.Stat(filepath.Join(currentDir, ".git")); err != nil {
				// Check if this is a worktree by looking for .git file
				gitFile := filepath.Join(currentDir, ".git")
				content, err := os.ReadFile(gitFile)
				if err != nil || !strings.Contains(string(content), "gitdir:") {
					return fmt.Errorf("not in a git repository or worktree")
				}
			}

			// Get current branch name
			currentBranch := strings.TrimSpace(runCommand("git", "rev-parse", "--abbrev-ref", "HEAD"))
			targetBranch := cmd.String("target")
			useRebase := cmd.Bool("rebase")
			shouldPush := cmd.Bool("push")

			// Fetch the latest changes without updating refs
			fmt.Printf("Fetching latest changes from remote...\n")
			runCommand("git", "fetch", "origin")

			// Instead of directly fetching into the target branch, update the merge base
			fmt.Printf("Updating %s branch reference...\n", targetBranch)
			// remoteRef := fmt.Sprintf("refs/remotes/origin/%s", targetBranch)

			// First, fetch the remote branch state
			runCommand("git", "fetch", "origin", targetBranch)

			if useRebase {
				fmt.Printf("Rebasing %s onto origin/%s...\n", currentBranch, targetBranch)
				output := runCommand("git", "rebase", fmt.Sprintf("origin/%s", targetBranch))
				if strings.Contains(strings.ToLower(output), "conflict") {
					return fmt.Errorf("rebase failed due to conflicts. Please resolve conflicts and continue the rebase manually")
				}
			} else {
				fmt.Printf("Merging origin/%s into %s...\n", targetBranch, currentBranch)
				output := runCommand("git", "merge", fmt.Sprintf("origin/%s", targetBranch))
				if strings.Contains(strings.ToLower(output), "conflict") {
					return fmt.Errorf("merge failed due to conflicts. Please resolve conflicts and commit the changes manually")
				}
			}

			if shouldPush {
				fmt.Printf("Pushing changes to remote...\n")
				if useRebase {
					runCommand("git", "push", "--force-with-lease", "origin", currentBranch)
				} else {
					runCommand("git", "push", "origin", currentBranch)
				}
			}

			fmt.Printf("Successfully updated %s with %s\n", currentBranch, targetBranch)
			return nil
		},
	}
}
