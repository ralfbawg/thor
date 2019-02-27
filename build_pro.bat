set GOARCH=amd64

set GOOS=linux

go build src/main/main.go "-w -s" -o bin