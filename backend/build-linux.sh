#!/bin/bash

# Load .env file
if [ -f ".env" ]; then
  export $(grep -v '^#' .env | xargs)
fi

OUTPUT="quicksaveService"

# Cross-compile for Linux (static binary)
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Build
go build -ldflags="-s -w -X 'main.clientID=${IGDB_API_KEY}' -X 'main.clientSecret=${IGDB_SECRET_KEY}'" -o "$OUTPUT"

echo "Linux static build complete: $OUTPUT"
