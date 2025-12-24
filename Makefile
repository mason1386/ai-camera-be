.PHONY: run-api run-worker build test lint docker-up docker-down gen-proto

run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

test:
	go test ./...

lint:
	golangci-lint run

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

gen-proto:
	sh scripts/gen_proto.sh

seed-db:
	sh scripts/seed_db.sh
