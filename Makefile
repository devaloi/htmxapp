.PHONY: build run test lint clean fmt

BINARY := htmxapp
CMD := ./cmd/htmxapp

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

test:
	go test -race -count=1 ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

clean:
	rm -f $(BINARY)
	go clean -testcache
