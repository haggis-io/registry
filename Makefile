VERSION=$(shell cat ./VERSION)
COMMIT=$(shell git rev-parse --short HEAD)
LATEST_TAG=$(shell git tag -l | head -n 1)

export VERSION COMMIT LATEST_TAG
.PHONY: test

test:
	@echo "=> Running tests"
	./hack/run-tests.sh

build:
	./hack/cross-platform-build.sh

verify:
	./hack/verify-version.sh

up: build
	draft up

generate:
	@echo "=> generating stubs"
	protoc -I ${PWD}/pkg/proto --proto_path=${PWD}/pkg/proto/ ${PWD}/pkg/proto/*.proto --go_out=plugins=grpc:${PWD}/pkg/proto
	@echo "=> injecting tags"
	protoc-go-inject-tag -input=./pkg/proto/entity.pb.go