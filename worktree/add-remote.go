package worktree

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

// AddRemoteWorktreeCommand creates a CLI command that:
// 1. Retrieves all remote branches.
// 2. Uses FZF to let you choose one.
// 3. Adds that branch as a worktree (with a sanitized directory name).
func AddRemoteWorktreeCommand() *cli.Command {
	return &cli.Command{
		Name:  "add-remote",
		Usage: "Select a remote branch via FZF and add it as a new worktree",
		Action: func(ctx context.Context, cliCmd *cli.Command) error {
			// Retrieve the remote branches.
			remoteOutput := runCommand("git", "branch", "-r")
			if remoteOutput == "" {
				return fmt.Errorf("no remote branches found")
			}

			// Split output into lines and filter out unwanted entries.
			lines := strings.Split(remoteOutput, "\n")
			var branches []string
			for _, line := range lines {
				branch := strings.TrimSpace(line)
				if branch == "" || strings.Contains(branch, "HEAD") {
					continue
				}
				branches = append(branches, branch)
			}
			if len(branches) == 0 {
				return fmt.Errorf("no valid remote branches found")
			}

			// Prepare the branch list for fzf.
			branchList := strings.Join(branches, "\n")

			// Start fzf and send the branch list via its stdin.
			fzf := exec.Command("fzf")
			stdin, err := fzf.StdinPipe()
			if err != nil {
				return fmt.Errorf("failed to get stdin pipe for fzf: %v", err)
			}
			go func() {
				defer stdin.Close()
				io.WriteString(stdin, branchList)
			}()

			// Capture the selected branch.
			selectedBytes, err := fzf.Output()
			if err != nil {
				return fmt.Errorf("fzf command failed or was canceled: %v", err)
			}
			selected := strings.TrimSpace(string(selectedBytes))
			if selected == "" {
				return fmt.Errorf("no branch selected")
			}

			// Remove the remote prefix ("origin/") to form the local branch name.
			localBranch := strings.TrimPrefix(selected, "origin/")

			// Get the repository's top-level directory.
			repoRoot := runCommand("git", "rev-parse", "--show-toplevel")
			if repoRoot == "" {
				return fmt.Errorf("failed to determine repository top-level directory")
			}

			// Place the new worktree as a sibling of the repo's top-level directory.
			parentDir := filepath.Dir(repoRoot)
			// Sanitize the branch name for a directory (replace "/" with "-").
			dirName := strings.ReplaceAll(localBranch, "/", "-")
			worktreePath := filepath.Join(parentDir, dirName)

			// Check if the target directory already exists.
			if _, err := os.Stat(worktreePath); err == nil {
				return fmt.Errorf("directory %s already exists", worktreePath)
			}

			fmt.Printf("Creating worktree at %s\n", worktreePath)
			// Create the worktree:
			// - Create a new local branch named `localBranch` that tracks the remote branch.
			// - Use the remote branch reference (e.g. "origin/feature/foo") as the start point.
			worktreeCmd := []string{
				"worktree", "add", worktreePath,
				"-b", localBranch,
				selected, // selected remote branch (e.g., "origin/feature/foo")
			}
			if output := runCommand("git", worktreeCmd...); output != "" {
				fmt.Println(output)
			}

			fmt.Printf("Successfully created worktree for branch '%s'\n", localBranch)
			return nil
		},
	}
}
