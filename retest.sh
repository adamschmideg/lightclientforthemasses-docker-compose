#!/bin/sh

go run faucet.go --template ./faucet.html &
sleep 1
PID=$(lsof -ti :8088)
go test -v
kill -9 $PID