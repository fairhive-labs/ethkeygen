.PHONY: run run-cli build build-cli clean test cover vet fmt tidy

run: clean build
	./bin/api
run-cli: clean build-cli
	./bin/ethkeygen $(COMMAND_ARGS)
build: clean
	go build -o bin/api -v api/main.go
build-cli: clean
	go build -o bin/ethkeygen -v ./cmd/main.go
clean:
	rm -rf ./bin
test:
	go clean -testcache
	go test -race -v ./...
cover:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
vet:
	go vet ./...
fmt:
	gofmt -l -w .
tidy:
	go mod tidy
