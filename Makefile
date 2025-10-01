BINARY_NAME=gonotes
MAIN=cmd/gonotes/main.go
MIGRATIONS_DIR=./migrations
STORAGE_PATH=sqlite://storage/storage.db

.PHONY: run build test clean lint fmt migrate-create migrate-up migrate-down migrate-reset migrate-version

run:
	go run $(MAIN)

build:
	go build -o $(BINARY_NAME) $(MAIN)

test:
	go test ./... -v

clean:
	go clean
	rm -f $(BINARY_NAME)

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database $(STORAGE_PATH) up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database $(STORAGE_PATH) down

migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database $(STORAGE_PATH) down

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database $(STORAGE_PATH) version
