#!/bin/bash

curl --location --request POST 'http://localhost:8080/v1/insert' \
--header 'Content-Type: application/json' \
--data-raw '{
    "dim": [{
        "key": "device",
        "val": "mobile"
},
{
    "key": "country",
    "val": "IN"
}],
"metrics": [{
    "key": "webreq",
    "val": 70
},
{
    "key": "timespent",
    "val": 30
}]
}'
