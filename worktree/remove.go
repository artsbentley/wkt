package worktree

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
)

func RemoveCommand() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove and cleanup worktree (selection via fzf)",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Get list of worktrees.
			output, err := exec.Command("git", "worktree", "list").Output()
			if err != nil {
				return err
			}

			// Launch fzf and pass the worktree list to it.
			fzfCmd := exec.Command("fzf")
			fzfCmd.Stdin = bytes.NewReader(output)
			selected, err := fzfCmd.Output()
			if err != nil {
				return err
			}

			// Convert the output to a string.
			line := string(selected)

			// The first field is the worktree path.
			fields := strings.Fields(line)
			if len(fields) == 0 {
				return nil // or return an error if no valid selection was made.
			}
			worktreePath := fields[0]

			// Attempt to extract the branch name from text in square brackets.
			re := regexp.MustCompile(`\[(.*?)\]`)
			matches := re.FindStringSubmatch(line)
			var branchName string
			if len(matches) >= 2 {
				branchName = matches[1]
			} else {
				// Fallback: use the worktree path as branch name.
				branchName = worktreePath
			}

			// Remove the worktree using the actual worktree path.
			runCommand("git", "worktree", "remove", "-f", worktreePath)

			// Prune the worktree list.
			runCommand("git", "worktree", "prune")

			// Remove branch from origin using the branch name.
			runCommand("git", "push", "origin", "--delete", branchName)

			// Remove local branch if prune didn't work.
			runCommand("git", "branch", "-d", branchName)

			// sync up with remote
			runCommand("git", "fetch", "origin", "--prune")

			return nil
		},
	}
}
