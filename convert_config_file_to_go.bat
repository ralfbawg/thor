set GOPATH=%~dp0

::go-bindata -prefix /src/config/file/ src/config/file/
go-bindata -o src/config/file.go  -pkg=config src/config/file/...