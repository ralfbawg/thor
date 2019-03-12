set GOPATH=%~dp0

set GOARCH=amd64

set GOOS=linux

go build  -o bin/thor src/main/main.go