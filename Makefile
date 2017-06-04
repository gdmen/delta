GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

build: build_server

build_server: src/server.go
	$(GOBUILD) -o ./bin/server $^

run: build
	./bin/server
