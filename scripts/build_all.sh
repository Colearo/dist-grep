#!/bin/bash
cd `dirname $0`
pwd
go build -o ../server/server ../server/server.go
go build -o ../client/client ../client/client.go
