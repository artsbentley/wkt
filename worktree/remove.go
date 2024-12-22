package worktree

import (
	"context"

	"github.com/urfave/cli/v3"
)

func RemoveCommand() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove and cleanup worktree",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			worktreeName := cmd.Args().First()

			// Configure the remote fetch
			runCommand("git", "worktree", "remove", "-f", worktreeName)

			// Remove local branch and its tracking branch
			runCommand("git", "worktree", "prune")

			// Remove branch from origin
			runCommand("git", "push", "origin", "--delete", worktreeName)

			// Remove local branch if prune didnt work
			runCommand("git", "branch", "-d", worktreeName)

			return nil
		},

		// Flags: []cli.Flag{
		// 	&cli.BoolFlag{
		// 		Name:    "tree",
		// 		Usage:   "create worktrees for all current branches",
		// 		Aliases: []string{"t"},
		// 		Value:   false,
		// 	},
		// },
	}
}
