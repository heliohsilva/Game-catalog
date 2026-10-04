.PHONY: all help build clean dev fe-dev fe-build api-build docker-build docker-up docker-down

all: fe-build

help:
	@echo "Game Catalog Makefile"
	@echo "Available commands:"
	@echo "  make              - Build frontend (default)"
	@echo "  make dev          - Run frontend development server"
	@echo "  make fe-build     - Build frontend application"
	@echo "  make api-build    - Build backend Golang API"
	@echo "  make clean        - Clean build artifacts (.next, dist, bin, cache)"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up    - Run application stack via docker-compose"
	@echo "  make docker-down  - Stop docker-compose stack"

dev:
	cd fe && npm run dev

fe-build:
	cd fe && npm run build

api-build:
	cd api && mkdir -p bin && go build -o bin/server main.go

clean:
	@echo "Cleaning build artifacts..."
	rm -rf fe/.next fe/out fe/.turbo api/bin
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
