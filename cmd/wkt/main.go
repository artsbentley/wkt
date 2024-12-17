package main

import (
	"context"
	"log"
	"os"

	"github.com/artsbentley/wkt/worktree"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "wkt",
		Usage: "CLI for everything git worktrees",
		Commands: []*cli.Command{
			worktree.WorktreeCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
