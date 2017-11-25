GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

local: api_server
	./bin/api_server > api_server.log 2>&1 &
	cd src/ui && npm start

test: test_api

test_api: src/api
	$(GOTEST) ./$^

api_server: src/api_server.go
	$(GOBUILD) -o ./bin/api_server ./$^

api_server_pi: src/api_server.go
	env GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o ./bin/api_server_pi ./$^

ui:
	cd src/ui && npm run build; cd -

release: api_server_pi ui
	mkdir -p ./bin/release
	cp ./bin/api_server_pi ./bin/release
	cp conf.json ./bin/release
	cp -r release/* ./bin/release
	cp -r src/ui/build ./bin/release/ui_server

deploy: release
	scp -r ./bin/release/* pi@10.0.0.174:./delta
