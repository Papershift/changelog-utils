package commands

import (
	"fmt"
	"os"

	"github.com/Papershift/changelog-utils/etc"
	"github.com/urfave/cli/v2"
)

var InitCommand = cli.Command{
	Name:  "init",
	Usage: "Initial setup",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "github-username",
			Usage:    "Your github username",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "fullname",
			Usage:    "Your full name",
			Required: true,
		},
	},
	Action: initHandler,
}

func initHandler(ctx *cli.Context) error {
	conf := etc.Config{
		User: etc.UserConfig{
			GithubUsername: ctx.String("github-username"),
			Fullname:       ctx.String("fullname"),
		},
	}

	if err := makeUserDir(conf.User.GithubUsername); err != nil {
		return err
	}

	if err := makeGitignoreFile(); err != nil {
		return err
	}

	if err := etc.MakeConfigFile(conf); err != nil {
		return err
	}

	return nil
}

func makeUserDir(name string) error {
	path := fmt.Sprintf("changelogs/%s", name)
	err := os.MkdirAll(path, 0755)

	if err == nil {
		fmt.Printf(">> %s/ created\n", path)
	}

	return err
}

func makeGitignoreFile() error {
	path := "changelogs/.gitignore"
	contents := []byte(etc.CONFIG_FILENAME)
	err := os.WriteFile(path, contents, 0644)

	if err == nil {
		fmt.Printf(">> %s written\n", path)
	}

	return err
}
