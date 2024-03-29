# changelog-utils

Set of CLI commands to manage CHANGELOG.md.
- Easy to add changelog entries and make releases
- Each changelog entry is separate file, therefore no Git conflicts
- Changelog entries are in human and machine readable [HCL](https://github.com/hashicorp/hcl) format

## Try
Download latest binary from Releases page and put it in empty directory.
Optionally rename the binary to `ch` for convenience.
Create new or move existing `CHANGELOG.md` file in the same directory. 
Then run commands:
```shell
# downloaded binary might not be executable, this makes it executable:
$ chmod +x ch

# this will create changelogs/ dir with local config file in it
$ ./ch init --github-username your-github-username --fullname "Your Name"

# this adds new changelog entry for "Added" section. Asana URL has to be close to 
# real Asana URL, because it is parsed when releasing to get last 4 digits. You can
# repeat similar command multiple times and instead of `added`, you can put `changed`,
# `fixed`, `removed`, `security`, `deprecated`
$ ./ch added --title "Example Asana task" --url https://app.asana.com/0/0/123456789/f

# Jira issue URLs are also supported
$ ./ch changed --title "Example Jira issue" --url "https://example.atlassian.net/jira/software/projects/PROJECT/boards/1?selectedIssue=PROJECT-123"
$ ./ch fixed --title "Another Jira issue" --url https://example.atlassian.net/browse/PROJECT-456

# following command will go through all changelog entries in changelogs/ subdirs, inserts
# them in CHANGELOG.md and removes entries from subdirs
$ ./ch release --version version-number
```
Once you are confident about usage, you can add the binary to existing project, where you
have `CHANGELOG.md` file and start using it.  

## Commands
### help
Prints help

### version
Prints current changelog-utils version

### init
**Params:**
- `--github-username` - Your Github username
- `--fullname` - Your full name  
Creates `changelogs/` directory with config file with the info you provided.
It also adds `.gitignore` to exclude config file from Git.

### added, changed, fixed, removed, security, deprecated
**Aliases:** `add`, `change`, `fix`, `remove`, `secure`, `deprecate`  
**Params:**
- `--title`, `-t` - Changelog entry title
- `--url`, `-u` - Jira issue or Asana task URL  
Adds new changelog entry under a section by command name (`added` => Added, etc.).

### release
**Params:**
- `--version`, `-v`  
Inserts all changelog entries into `CHANGELOG.md` under given version number and
then deletes the entries.  
**IMPORTANT**: This command searches for latest existing release in `CHANGELOG.md`
and the format of releaes heading is expected to be as follows: `## [00.0] - YYYY-MM-DD`

## TODO
- Unit tests
- Document code
- Configurable CHANGELOG.md file name
- Support for multiple task URLs
- Support for multiple authors
- Task URL format validator
- Version input validation
- Interactive commands
