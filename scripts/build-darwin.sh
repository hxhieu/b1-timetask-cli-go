#!/bin/bash
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o build/bin/b1-timetask-cli-go-darwin
