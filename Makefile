gofmt:
	go fmt ./...

test: gofmt
	go test ./...