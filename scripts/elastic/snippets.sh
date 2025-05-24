#/bin/sh

curl -X PUT "localhost:9200/snippets" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "name": {
        "type": "completion"
      },
      "description": {
        "type": "text"
      },
      "body": {
        "type": "text"
      },
      "userId": {
        "type": "keyword"
      }
    }
  }
}
'

# curl -H "Content-Type: application/json" http://localhost:9200/notes/_search?pretty=true -d '{ "query": { "match": { "name": "файл" } } }'