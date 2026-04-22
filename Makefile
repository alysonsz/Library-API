.PHONY: build run test vet clean docker-up docker-down

build:
	go build -o bin/gobook ./cmd/gobook

run:
	go run ./cmd/gobook

run-postgres:
	DB_DRIVER=postgres DB_URL=postgres://library:library@localhost:5432/library?sslmode=disable go run ./cmd/gobook

test:
	go test ./... -v

vet:
	go vet ./...

tidy:
	go mod tidy

swag:
	swag init -g cmd/gobook/main.go

clean:
	go clean
	rm -f bin/gobook
	rm -f *.db

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

cli-search:
	go run ./cmd/gobook search "$(term)"

cli-simulate:
	go run ./cmd/gobook simulate $(ids)
