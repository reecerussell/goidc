#!/bin/bash

set -e

echo "Building $NAME v$VERSION..."
echo "Working Directory: $WORKING_DIRECTORY"

cd "$WORKING_DIRECTORY"

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=false

echo -e "Env:\n-GOOS=$GOOS\n-GOARCH=$GOARCH\n-CGO_ENABLED=$CGO_ENABLED"

echo "Installing modules..."
go mod download

echo "Building..."
go build -o main main.go

echo "Built successfully!"

echo "Zipping..."

zip_name="build.zip"
zip $zip_name main

echo "Zipped: $zip_name"
