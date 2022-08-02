package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Papershift/changelog-utils/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "changelog-utils",
		Usage: "add changelog and release entries",

		Commands: []*cli.Command{
			&commands.InitCommand,
			&commands.ReleaseCommand,
			&commands.VersionCommand,

			&commands.AddedCommand,
			&commands.ChangedCommand,
			&commands.FixedCommand,
			&commands.RemovedCommand,
			&commands.SecurityCommand,
			&commands.DeprecatedCommand,
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
