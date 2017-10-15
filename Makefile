GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

build: build_server

build_server: src/server.go
	$(GOBUILD) -o ./bin/server ./$^

test: test_api

test_api: src/api
	$(GOTEST) ./$^

run_api: build
	./bin/server > server.log 2>&1

run_ui:
	cd src/ui && npm start
