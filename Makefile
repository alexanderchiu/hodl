VERSION := $(shell cat ./VERSION)
all: install

install:
	go install -v ./...

test:
	go test ./... -v

fmt:
	go fmt ./...

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

build:
	mkdir -p bin
	go build -o bin ./...

.PHONY: install test fmt release