#!/bin/bash

ARCH=$1

[ "$ARCH" != "amd64" ] && [ "$ARCH" != "arm64" ] && echo "Usage: $0 amd64|arm64" && exit 1

cd /app

NOW=$(date +'%Y-%m-%d_%H:%M:%S') 
VERSION=`cat version`

export GGO_ENABLED=1


[ "$ARCH" == "amd64" ] &&  CGO_CFLAGS="-O -D__BLST_PORTABLE__" GOARCH=amd64 GOOS=linux go build -ldflags "-X main.buildTime=${NOW}_UTC -X main.version=${VERSION}" -o rocketreport-${ARCH} cmd/rocketreport/main.go

[ "$ARCH" == "arm64" ] &&  CGO_CFLAGS="-O -D__BLST_PORTABLE__" GOARCH=arm64 GOOS=linux CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-cpp go build -o rocketpool-${ARCH} cmd/rocketreport/main.go
