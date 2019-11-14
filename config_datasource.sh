#!/bin/sh

curl -u'admin:foobar' \
 -X POST \
 -H "Content-Type: application/json" \
 -d@datasource.json \
 http://localhost:3000/api/datasources
