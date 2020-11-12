gofmt:
	go fmt ./...

test: gofmt
	go test ./... -v

run: test
	docker-compose up --force-recreate --build app


