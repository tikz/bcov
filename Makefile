VERSION = $(shell git describe --tags)
COMMIT = $(shell git rev-parse --short HEAD)

all: build

build: dep
	npm version $(VERSION) --prefix web/
	npm run build --prefix web/
	go build -ldflags="-X 'main.Version=$(VERSION)' -X 'main.CommitHash=$(COMMIT)'"

test: dep
	go test ./... -v

lint: dep
	golint

dep:
	go mod tidy