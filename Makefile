gofmt:
	go fmt ./...

test: gofmt
	go test ./... -v

run: test
	docker-compose up -d --force-recreate --build app

testsuite: 
	sh scripts/curl_insert.sh
	sh scripts/curl_query.sh 

stop:
	docker-compose down

clean:
	rm -rf curl_insert curl_query


