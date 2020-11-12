#!/bin/bash

curl --location --request POST 'http://localhost:8080/v1/query' \
--header 'Content-Type: application/json' \
--data-raw '{
	"dim": [{
	"key": "country",
	"val": "IN"
}]
}'
