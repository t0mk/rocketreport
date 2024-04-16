#!/bin/bash

cd /app

CGO_CFLAGS="-O -D__BLST_PORTABLE__" GOARCH=amd64 GOOS=linux go build -o rocketreport-amd64 cmd/rocketreport/main.go


