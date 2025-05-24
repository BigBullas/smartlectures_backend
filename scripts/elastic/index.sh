#/bin/sh

curl -X PUT "localhost:9200/notes" -H 'Content-Type: application/json' -d "$(cat ./configs/elastic/index.json)"

