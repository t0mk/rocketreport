#!/bin/bash

cd /app

NOW=$(date +'%Y-%m-%d_%H:%M:%S') 
VERSION=`cat version`

CGO_ENABLED=1 CGO_CFLAGS="-O -D__BLST_PORTABLE__" GOARCH=amd64 GOOS=linux go build -ldflags "-X main.buildTime=${NOW}_UTC -X main.version=${VERSION}" -o rocketreport-amd64 cmd/rocketreport/main.go
