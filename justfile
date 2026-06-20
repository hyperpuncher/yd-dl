build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/yd-dl-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=arm64 go build -o bin/yd-dl-darwin-arm64 .

build-mac-x64:
	GOOS=darwin GOARCH=amd64 go build -o bin/yd-dl-darwin-amd64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/yd-dl-windows-amd64.exe .

fmt:
	gofmt -w .

vet: fmt
	go vet ./...

build-all: vet build-linux build-mac build-mac-x64 build-windows

clean:
	trash bin
