BINARY_NAME=app

run:
	go run cmd/gonotes/main.go

build:
	go build -o $(BINARY_NAME) main.go

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY_NAME)
