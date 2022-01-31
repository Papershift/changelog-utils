package main

import (
	"fmt"
	"log"
	"os"

	"github.com/azamat-sharapov/changelog-utils/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "changelog-utils",
		Usage: "add changelog and release entries",

		Commands: []*cli.Command{
			&commands.InitCommand,

			&commands.AddedCommand,
			&commands.ChangedCommand,
			&commands.FixedCommand,
			&commands.RemovedCommand,
			&commands.SecurityCommand,
			&commands.DeprecatedCommand,

			&commands.ReleaseCommand,
		},

		Before: func(ctx *cli.Context) error {
			// add empty line
			fmt.Println()
			return nil
		},
		After: func(ctx *cli.Context) error {
			fmt.Println()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
