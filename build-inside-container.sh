#!/bin/bash

cd /app

CGO_ENABLED=1 CGO_CFLAGS="-O -D__BLST_PORTABLE__" GOARCH=amd64 GOOS=linux go build -o rocketreport cmd/rocketreport/main.go