build:
	GOOS=darwin GOARCH=amd64 go build -o bin/ch-amd64-darwin
	GOOS=linux GOARCH=amd64 go build -o bin/ch-amd64-linux
	GOOS=windows GOARCH=amd64 go build -o bin/ch-amd64.exe
