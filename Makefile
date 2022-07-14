VERSION = $(shell git describe --tags)
COMMIT = $(shell git rev-parse --short HEAD)

all: build

build: export CGO_ENABLED = 0
build: dep
	npm version $(VERSION)-$(COMMIT) --allow-same-version --prefix web/
	npm install --legacy-peer-deps --prefix web/
	npm run build --prefix web/
	go mod tidy
	go build -ldflags="-X 'main.Version=$(VERSION)' -X 'main.CommitHash=$(COMMIT)'"

test: dep
	go test ./... -v

lint: dep
	golint

dep:
	go mod tidy