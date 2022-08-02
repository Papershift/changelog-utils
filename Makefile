VERSION="$$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)"
LDFLAGS="-X 'github.com/Papershift/changelog-utils/commands.Version=${VERSION}'"

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags=${LDFLAGS} -o bin/ch-amd64-darwin
	GOOS=linux GOARCH=amd64 go build -ldflags=${LDFLAGS} -o bin/ch-amd64-linux
	GOOS=windows GOARCH=amd64 go build -ldflags=${LDFLAGS} -o bin/ch-amd64.exe
