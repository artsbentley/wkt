package worktree

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func WorktreeCommand() *cli.Command {
	return &cli.Command{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}
}
