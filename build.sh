#!/bin/sh

version="v2.0.0"
echo $version

# Enable parallel builds
build() {
  CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -trimpath -ldflags "-s -w -X 'main.VERSION=$version'" -o bin/$1_$2/ cmd/pzrpc/pzrpc.go &
  CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -trimpath -ldflags "-s -w -X 'main.VERSION=$version'" -o bin/$1_$2/ cmd/pzrps/pzrps.go &
}

# Build for Linux, Windows, macOS and add ARM architectures
build linux amd64
build linux arm64
build windows amd64
build darwin amd64
build darwin arm64

wait  # Wait for all background processes to finish

echo "Build completed!"