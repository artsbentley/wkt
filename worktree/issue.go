package worktree

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
)

func IssueCommand() *cli.Command {
	return &cli.Command{
		Name:  "issue",
		Usage: "Create a worktree from a selected GitHub issue",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Step 1: Get GitHub issues using gh CLI
			ghCmd := exec.Command("gh", "issue", "list", "--limit", "100", "--state", "open", "--json", "number,title", "--template", "{{range .}}{{.number}}: {{.title}}\n{{end}}")
			output, err := ghCmd.Output()
			if err != nil {
				return fmt.Errorf("failed to fetch issues: %w", err)
			}

			// Step 2: Let user pick one with fzf
			fzf := exec.Command("fzf")
			fzf.Stdin = bytes.NewReader(output)
			selected, err := fzf.Output()
			if err != nil {
				return fmt.Errorf("fzf selection failed: %w", err)
			}

			// Step 3: Parse issue number and title
			line := strings.TrimSpace(string(selected))
			matches := regexp.MustCompile(`^(\d+):\s+(.+)$`).FindStringSubmatch(line)
			if len(matches) != 3 {
				return fmt.Errorf("could not parse selected issue line: %q", line)
			}
			issueNum := matches[1]
			title := matches[2]

			// Sanitize for branch/directory name
			safeTitle := sanitizeBranchPart(title)
			branch := fmt.Sprintf("%s-%s", issueNum, safeTitle)

			// Step 4: Determine worktree path.
			repoRoot := runCommand("git", "rev-parse", "--show-toplevel")
			if repoRoot == "" {
				return fmt.Errorf("failed to get repository top-level directory")
			}
			parentDir := filepath.Dir(repoRoot)
			worktreePath := filepath.Join(parentDir, branch)

			if _, err := os.Stat(worktreePath); err == nil {
				return fmt.Errorf("worktree directory %s already exists", worktreePath)
			}

			fmt.Printf("Creating worktree at %s with branch %s based on main\n", worktreePath, branch)

			// Step 5: Create the worktree and push branch
			runCommand("git", "worktree", "add", worktreePath, "-b", branch, "main")
			if err := os.Chdir(worktreePath); err != nil {
				return fmt.Errorf("failed to change to directory %s: %v", worktreePath, err)
			}
			runCommand("git", "push", "-u", "origin", branch)

			fmt.Printf("Successfully created worktree for issue #%s\n", issueNum)
			return nil
		},
	}
}

// sanitizeBranchPart makes a string safe for use in branch and folder names.
func sanitizeBranchPart(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, ",", "")
	return s
}
