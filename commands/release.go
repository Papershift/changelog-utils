package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const RELEASE_VERSION_PATTERN = `##\s+\[\d+.?\d+\]\s+-\s+\d{4}-\d{2}-\d{2}\s?`
const RELEASE_VERSION_FORMAT = "## [%s] - %s\n"
const NEW_LINE = '\n'
const (
	SECTION_ADDED      = "added"
	SECTION_CHANGED    = "changed"
	SECTION_FIXED      = "fixed"
	SECTION_REMOVED    = "removed"
	SECTION_SECURITY   = "security"
	SECTION_DEPRECATED = "deprecated"
)

var ALL_SECTIONS = [...]string{SECTION_ADDED, SECTION_CHANGED, SECTION_FIXED, SECTION_REMOVED, SECTION_SECURITY, SECTION_DEPRECATED}

var ReleaseCommand = cli.Command{
	Name:  "release",
	Usage: "Make a release",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "version",
			Aliases:  []string{"v"},
			Usage:    "Releaes version number",
			Required: true,
		},
	},
	Action: releaseHandler,
}

func releaseHandler(ctx *cli.Context) error {
	tmpfile, err := ioutil.TempFile("changelogs", ".ch.tmp.")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	beginning, end, err := getChangelogParts()
	if err != nil {
		return err
	}

	_, err = tmpfile.Write(beginning)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	version := fmt.Sprintf(RELEASE_VERSION_FORMAT, ctx.String("version"), currentTime.Format("2006-01-02"))

	usersEntries, err := getUsersEntriesPaths()
	if err != nil {
		return err
	}

	if _, err = tmpfile.Write([]byte(version)); err != nil {
		return err
	}

	content, err := makeReleaseContent(usersEntries)
	if err != nil {
		return err
	}
	if _, err = tmpfile.Write(content); err != nil {
		return err
	}

	if _, err = tmpfile.WriteString(string(NEW_LINE)); err != nil {
		return err
	}

	if _, err = tmpfile.Write(end); err != nil {
		return err
	}

	if err = tmpfile.Sync(); err != nil {
		return err
	}

	if err = os.Rename(tmpfile.Name(), "CHANGELOG.md"); err != nil {
		return err
	}

	fmt.Printf(">> release %s written to %s", ctx.String("version"), "CHANGELOG.md")
	fmt.Println()

	if err = deleteUsersEntries(usersEntries); err != nil {
		return err
	}

	return nil
}

func makeReleaseContent(usersEntries []string) ([]byte, error) {
	sections := map[string]*bytes.Buffer{
		SECTION_ADDED:      {},
		SECTION_CHANGED:    {},
		SECTION_FIXED:      {},
		SECTION_REMOVED:    {},
		SECTION_SECURITY:   {},
		SECTION_DEPRECATED: {},
	}

	for _, userEntry := range usersEntries {
		entry, err := ReadEntryFile(userEntry)
		if err != nil {
			return nil, err
		}

		if _, ok := sections[entry.Section]; !ok {
			return nil, errors.New(fmt.Sprintf("err: unknown section %s", entry.Section))
		}

		buff := sections[entry.Section]
		buff.Write(makeReleaseEntry(entry))
	}

	var content bytes.Buffer

	for _, section := range ALL_SECTIONS {
		entries := sections[section]
		if entries.Len() == 0 {
			continue
		}

		title := strings.Title(section)

		content.WriteString(fmt.Sprintf("### %s", title))
		content.WriteRune(NEW_LINE)
		content.Write(entries.Bytes())

		fmt.Printf(">> found entries for %s", title)
		fmt.Println()
	}

	return content.Bytes(), nil
}

func getChangelogParts() ([]byte, []byte, error) {
	var beginning bytes.Buffer
	var end bytes.Buffer

	file, err := os.OpenFile("CHANGELOG.md", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lastReleaseFound := false

	for scanner.Scan() {
		currentLine := scanner.Bytes()

		if !lastReleaseFound {
			if matched, err := regexp.Match(RELEASE_VERSION_PATTERN, currentLine); err != nil {
				return nil, nil, err
			} else if matched {
				fmt.Printf(">> last release found: %s", currentLine[3:])
				fmt.Println()

				lastReleaseFound = true
				end.Write(currentLine)
				end.WriteRune(NEW_LINE)
			} else {

				beginning.Write(currentLine)
				beginning.WriteRune(NEW_LINE)
			}
		} else {
			end.Write(currentLine)
			end.WriteRune(NEW_LINE)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return beginning.Bytes(), end.Bytes(), nil
}

func makeReleaseEntry(entry *Entry) []byte {
	var str bytes.Buffer

	asanaShortNumb := getAsanaUrlShortNumb(entry.Asana)

	firstLine := fmt.Sprintf("- %s  ", entry.Title)
	secondLine := fmt.Sprintf("  Task: [%s](%s) | Author: [%s](%s)", asanaShortNumb, entry.Asana, entry.Author, entry.Github)

	str.WriteString(firstLine)
	str.WriteRune(NEW_LINE)
	str.WriteString(secondLine)
	str.WriteRune(NEW_LINE)

	return str.Bytes()
}

func deleteUsersEntries(entriesPaths []string) error {
	for _, path := range entriesPaths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func getUsersEntriesPaths() ([]string, error) {
	var entries []string

	changelogsPath := "changelogs"
	changelogsDir, err := os.ReadDir(changelogsPath)
	if err != nil {
		return nil, err
	}

	for _, changelogDirItem := range changelogsDir {
		if !changelogDirItem.IsDir() {
			continue
		}

		userDirPath := fmt.Sprintf("%s/%s", changelogsPath, changelogDirItem.Name())
		userDir, err := os.ReadDir(userDirPath)
		if err != nil {
			return nil, err
		}

		for _, userDirEntry := range userDir {
			if !userDirEntry.IsDir() || strings.HasSuffix(userDirEntry.Name(), ".hcl") {
				entries = append(entries, fmt.Sprintf("%s/%s", userDirPath, userDirEntry.Name()))
			}
		}
	}

	return entries, nil
}

func getAsanaUrlShortNumb(url string) string {
	split := strings.Split(url, "/")

	for i := len(split) - 1; i >= 0; i-- {
		str := split[i]

		if len(str) < 4 {
			continue
		}

		if _, err := strconv.Atoi(str); err == nil {
			return str[len(str)-4:]
		}
	}

	return ""
}
