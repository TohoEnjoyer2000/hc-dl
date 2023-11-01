default:
	mkdir -p build
	GOOS=darwin GOARCH=arm64 go build -o build/hc-dl-macos-arm64 cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o build/hc-dl-macos-amd64 cmd/main.go
	GOOS=linux GOARCH=amd64 go build -o build/hc-dl-linux-amd64 cmd/main.go
	GOOS=linux GOARCH=arm64 go build -o build/hc-dl-linux-arm64 cmd/main.go
	GOOS=linux GOARCH=arm go build -o build/hc-dl-linux-arm cmd/main.go
	GOOS=windows GOARCH=amd64 go build -o build/hc-dl.exe cmd/main.go