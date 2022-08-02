package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var VersionCommand = cli.Command{
	Name:   "version",
	Usage:  "Print changelog-utils version",
	Action: versionHandler,
}

var Version string = "dev"

func versionHandler(ctx *cli.Context) error {
	fmt.Println(Version)
	return nil
}
