#! /bin/bash

mkdir -p bin
GOOS=linux go build -o ./bin/upload ./upload.go
GOOS=linux go build -o ./bin/download ./download.go