package worktree

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"
)

// ANSI escape codes for color formatting
const (
	reset = "\033[0m"
	gray  = "\033[0;90m"
	red   = "\033[0;31m"
	green = "\033[0;32m"
	bold  = "\033[1m"
)

var (
	verbose          bool
	branch           string
	base             string
	prefix           string
	noCreateUpstream bool
	githubPrefix     bool
	worktree         string
)

func main() {
	rootFlagSet := flag.NewFlagSet("gitworktree", flag.ExitOnError)
	rootFlagSet.BoolVar(&verbose, "verbose", false, "Print script debug info")
	rootFlagSet.StringVar(&branch, "branch", "", "The branch to create")
	rootFlagSet.StringVar(&base, "base", "origin/main", "The branch to use as the base for the new worktree")
	rootFlagSet.BoolVar(&noCreateUpstream, "no-create-upstream", false, "Do not create an upstream branch")
	rootFlagSet.BoolVar(&githubPrefix, "github-prefix", false, "Use GitHub username as the prefix for the branch name")

	root := &ffcli.Command{
		Name:       "gitworktree",
		ShortUsage: "gitworktree <path> [flags]",
		ShortHelp:  "Script the creation of a new worktree.",
		FlagSet:    rootFlagSet,
		Exec:       run,
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing script arguments")
	}
	worktree = args[0]

	if githubPrefix {
		prefix = getGitConfig("github.user") + "/"
	}

	if branch == "" {
		branch = prefix + filepath.Base(worktree)
	}

	if branchExists(branch) {
		fmt.Printf("Branch %s already exists, adding worktree %s\n", branch, worktree)
		if err := runCommand(fmt.Sprintf("%sGenerating new worktree from existing branch: %s%s\n\n", bold, branch, reset), "git", "worktree", "add", worktree, branch); err != nil {
			return err
		}
	} else {
		fmt.Printf("Branch %s does not exist, creating new worktree %s based on %s\n", branch, worktree, base)
		if err := runCommand(fmt.Sprintf("%sGenerating new worktree: %s%s\n\n", bold, worktree, reset), "git", "worktree", "add", "-b", branch, worktree, base); err != nil {
			return err
		}
	}

	fmt.Printf("%sMoving into worktree: %s%s\n\n", gray, worktree, reset)
	if err := os.Chdir(worktree); err != nil {
		return err
	}
	if !noCreateUpstream {
		if err := updateRemote(branch); err != nil {
			return err
		}
	}
	fmt.Printf("%sSuccess.%s\n\n", green, reset)
	return nil
}

func getGitConfig(key string) string {
	out, err := execCommand("git", "config", "--get", key)
	if err != nil {
		log.Fatalf("failed to get git config %s: %v", key, err)
	}
	return strings.TrimSpace(out)
}

func branchExists(branch string) bool {
	out, err := execCommand("git", "branch", "--list", branch)
	if err != nil {
		log.Fatalf("failed to list git branches: %v\nOutput: %s", err, out)
	}
	return strings.TrimSpace(out) != ""
}

func updateRemote(branch string) error {
	out, err := execCommand("git", "ls-remote", "--heads", "origin", branch)
	if err != nil {
		fmt.Printf("failed to check remote branches: %v\nOutput: %s", err, out)
		return err
	}
	if strings.TrimSpace(out) == "" {
		fmt.Printf("%sBranch '%s' does not exist on remote. Creating.%s\n\n", gray, branch, reset)
		if err := runCommand(fmt.Sprintf("%sCreating remote branch %s%s\n\n", bold, branch, reset), "git", "push", "--set-upstream", "origin", branch); err != nil {
			return err
		}
	} else {
		fmt.Printf("%sBranch '%s' exists. Setting upstream.%s\n\n", gray, branch, reset)
		if err := runCommand(fmt.Sprintf("%sSetting upstream branch to 'origin/%s'%s\n\n", bold, branch, reset), "git", "branch", "--set-upstream-to=origin/"+branch); err != nil {
			return err
		}
	}
	return nil
}

func execCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("%sCommand '%s %s' failed with exit status %d: %s%s\n\n", red, command, strings.Join(args, " "), exitError.ExitCode(), string(exitError.Stderr), reset)
		}
		return "", err
	}
	return string(out), nil
}

func runCommand(message string, command string, args ...string) error {
	fmt.Print(message + " ")
	cmd := exec.Command(command, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("%sFAILED.%s\n\n", red, reset)
		return err
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	spinner := []rune{'⣾', '⣽', '⣻', '⢿', '⡿', '⣟', '⣯', '⣷'}
	i := 0

loop:
	for {
		select {
		case err := <-done:
			if err != nil {
				fmt.Printf("%sFAILED.%s\n\n", red, reset)
				return err
			}
			fmt.Printf("%sDone.%s\n\n", green, reset)
			break loop
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("\r%c ", spinner[i])
			i = (i + 1) % len(spinner)
		}
	}

	return nil
}
