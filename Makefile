all: build

build: dep
	go build -ldflags="-X 'main.Version=v1.0.0'"

test: dep
	go test ./... -v

lint: dep
	golint

dep:
	go mod tidy