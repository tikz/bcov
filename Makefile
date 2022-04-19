VERSION = $(shell git describe --tags)
COMMIT = $(shell git rev-parse --short HEAD)

all: build

build: dep
	go build -ldflags="-X 'main.Version=$(VERSION)' -X 'main.CommitHash=$(COMMIT)'"

test: dep
	go test ./... -v

lint: dep
	golint

dep:
	go mod tidy