package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Papershift/changelog-utils/etc"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/urfave/cli/v2"
	"github.com/zclconf/go-cty/cty"
)

type Entry struct {
	Section string `hcl:"section"`
	Title   string `hcl:"title"`
	Url     string `hcl:"url"`
	Author  string `hcl:"author"`
	Github  string `hcl:"github"`
}

var AddedCommand = makeEntryCommand("added", "add")
var ChangedCommand = makeEntryCommand("changed", "change")
var FixedCommand = makeEntryCommand("fixed", "fix")
var RemovedCommand = makeEntryCommand("removed", "remove")
var SecurityCommand = makeEntryCommand("security", "secure")
var DeprecatedCommand = makeEntryCommand("deprecated", "deprecate")

func makeEntryCommand(name string, alias string) cli.Command {
	return cli.Command{
		Name:    name,
		Aliases: []string{alias},
		Usage:   fmt.Sprintf("Add new changelog entry under %s section", strings.Title(name)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "Entry title",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "Jira issue or Asana task URL",
				Required: true,
			},
		},
		Action: entryHandler,
	}
}

func entryHandler(ctx *cli.Context) error {
	if !etc.ConfigExists() {
		return errors.New("config file missing. Run `init` command")
	}

	conf, err := etc.ReadConfig()
	if err != nil {
		return err
	}

	entry := Entry{
		Section: ctx.Command.Name,
		Title:   ctx.String("title"),
		Url:     ctx.String("url"),
		Author:  conf.User.Fullname,
		Github:  fmt.Sprintf("https://github.com/%s", conf.User.GithubUsername),
	}

	path := fmt.Sprintf("changelogs/%s/%d.hcl", conf.User.GithubUsername, time.Now().Unix())
	file := hclwrite.NewEmptyFile()
	body := file.Body()

	body.SetAttributeValue("section", cty.StringVal(entry.Section))
	body.SetAttributeValue("title", cty.StringVal(entry.Title))
	body.SetAttributeValue("url", cty.StringVal(entry.Url))
	body.SetAttributeValue("author", cty.StringVal(entry.Author))
	body.SetAttributeValue("github", cty.StringVal(entry.Github))

	err = os.WriteFile(path, file.Bytes(), 0644)

	if err == nil {
		fmt.Printf(">> %s written\n", path)
	}

	return nil
}

func ReadEntryFile(path string) (*Entry, error) {
	var entry Entry

	err := hclsimple.DecodeFile(path, nil, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
