run: clean build
	./bin/ethkeygen $(COMMAND_ARGS)
build: clean
	go build -o bin/ethkeygen -v ./cmd/main.go 
clean:
	rm -rf ./bin
test:
	go clean -testcache
	go test -v ./... 