#!/bin/bash
set -e

# Install Go if not present
if ! command -v go &> /dev/null; then
    curl -sL https://go.dev/dl/go1.22.0.linux-amd64.tar.gz | tar -C /usr/local -xzf -
    export PATH=$PATH:/usr/local/go/bin
fi

# Download Go dependencies
go mod download

# Build frontend
cd web
npm ci
npm run build

# Copy to public
cp -r dist/* ../public/ 2>/dev/null || cp -r build/* ../public/
