.PHONY: all help build clean clean-all dev fe-dev fe-build api-build api-dev api-test test docker-build docker-up docker-down

all: fe-build

help:
	@echo "Game Catalog Makefile"
	@echo "Available commands:"
	@echo "  make              - Build frontend (default)"
	@echo "  make dev          - Run frontend development server"
	@echo "  make fe-build     - Build frontend application"
	@echo "  make api-build    - Build backend Golang API"
	@echo "  make api-dev      - Run backend Golang API locally"
	@echo "  make api-test     - Run backend API unit and integration tests"
	@echo "  make test         - Run test suite"
	@echo "  make clean        - Clean build artifacts and database files"
	@echo "  make clean-all    - Deep clean including node_modules"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up    - Run application stack via docker-compose"
	@echo "  make docker-down  - Stop docker-compose stack"

dev:
	cd fe && npm run dev

fe-build:
	cd fe && npm run build

api-build:
	cd api && mkdir -p bin && go build -o bin/server main.go

api-dev:
	cd api && go run main.go

api-test:
	cd api && go test -v ./...

test: api-test

clean:
	@echo "Cleaning build artifacts..."
	rm -rf fe/.next fe/out fe/.turbo api/bin api/*.db api/*.db-shm api/*.db-wal
	@echo "Clean completed."

clean-all: clean
	@echo "Deep cleaning node_modules..."
	rm -rf fe/node_modules
	@echo "Deep clean completed."

docker-build:
	docker compose -f env/docker-compose.yml build

docker-up:
	docker compose -f env/docker-compose.yml up -d

docker-down:
	docker compose -f env/docker-compose.yml down
