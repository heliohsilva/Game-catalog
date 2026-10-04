.PHONY: all help build clean clean-all dev fe-dev fe-build fe-test api-build api-dev api-test test docker-build docker-up docker-down

all: docker-build docker-up

build: docker-build

help:
	@echo "Game Catalog Makefile"
	@echo "Available commands:"
	@echo "  make              - Build whole environment and start containers (docker-build + docker-up)"
	@echo "  make build        - Build all Docker images"
	@echo "  make test         - Run all test suites (frontend + backend)"
	@echo "  make fe-test      - Run frontend Vitest test suite"
	@echo "  make api-test     - Run backend API unit and integration tests"
	@echo "  make fe-build     - Build frontend Next.js production bundle locally"
	@echo "  make api-build    - Build backend Golang API binary locally"
	@echo "  make dev          - Run frontend development server"
	@echo "  make api-dev      - Run backend Golang API locally"
	@echo "  make clean        - Clean build artifacts and database files"
	@echo "  make clean-all    - Deep clean including node_modules"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up    - Run application stack via docker-compose"
	@echo "  make docker-down  - Stop docker-compose stack"

dev:
	cd fe && npm run dev

fe-dev:
	cd fe && npm run dev

fe-build:
	cd fe && npm run build

fe-test:
	cd fe && npm test

api-build:
	cd api && mkdir -p bin && go build -o bin/server main.go

api-dev:
	cd api && go run main.go

api-test:
	cd api && go test -v ./...

test: fe-test api-test

seed-db:
	cd api && go run cmd/seed/main.go -path ../data/game_catalog.sqlite
	cp data/game_catalog.sqlite game_catalog.sqlite

clean:
	@echo "Cleaning build artifacts..."
	rm -rf fe/.next fe/out fe/.turbo api/bin *.db *.db-shm *.db-wal *.sqlite-shm *.sqlite-wal api/*.db api/*.sqlite-shm api/*.sqlite-wal
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
