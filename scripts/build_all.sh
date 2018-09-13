#!/bin/bash
cd `dirname $0`
#pwd
go build -o ../src/grep-server/server ../src/grep-server/server.go
