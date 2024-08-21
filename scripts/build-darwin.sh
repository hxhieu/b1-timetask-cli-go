#!/bin/bash
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o .bin/timetask-darwin
