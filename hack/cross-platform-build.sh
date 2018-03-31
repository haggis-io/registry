#!/bin/bash
for GOOS in darwin linux windows; do
	for GOARCH in 386 amd64; do
        echo "Building binary for $GOOS on $GOARCH"
        GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/${GOOS}/${GOARCH}/registry -ldflags "-X main.commit=${COMMIT} -X main.version=${VERSION}" cmd/registry/main.go
    done
done
echo "Complete"