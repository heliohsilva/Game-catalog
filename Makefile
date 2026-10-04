.PHONY: all clean test

all:
	docker compose -f env/docker-compose.yml build
	docker compose -f env/docker-compose.yml up -d

clean:
	docker compose -f env/docker-compose.yml down
	rm -rf fe/.next fe/out fe/.turbo api/bin *.db *.db-shm *.db-wal *.sqlite-shm *.sqlite-wal

test:
	cd api && go test ./...
	cd fe && npm test
