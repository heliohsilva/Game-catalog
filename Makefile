.PHONY: all clean test

all:
	@if [ ! -f env/.env ]; then cp env/.env.example env/.env && echo "Created env/.env from env/.env.example"; fi
	docker compose -f env/docker-compose.yml build
	docker compose -f env/docker-compose.yml up -d

migrate:
	@if [ ! -f data/game_catalog.sqlite ]; then \
		cp ../game_catalog_bd/game_catalog.sql data/game_catalog.sql; \
		cd data && ./sql_mount; \
	fi

clean:
	docker compose -f env/docker-compose.yml down
	rm -rf fe/.next fe/out fe/.turbo api/bin *.db *.db-shm *.db-wal *.sqlite-shm *.sqlite-wal

test:
	cd api && go test ./...
	cd fe && npm test
