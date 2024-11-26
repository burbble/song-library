.PHONY: up down migrate postgres recreate-db build logs start reset-db restart-app

DC=docker compose
DB_USER=song_library_user
DB_PASS=song_library_password
DB_NAME=song_library_db
DB_HOST=localhost
DB_PORT=5432
MIGRATIONS_DIR=migrations
PROJECT_NAME=song-library
NETWORK=$(PROJECT_NAME)_song-network

up:
	$(DC) up

down:
	$(DC) down -v

build:
	$(DC) build

rebuild:
	$(DC) up --build

logs:
	@if [ "$(service)" = "" ]; then \
		$(DC) logs -f; \
	else \
		$(DC) logs -f $(service); \
	fi

postgres:
	$(DC) up -d postgres
	@until docker exec $$(docker ps -q -f name=postgres) pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		echo "Waiting for postgres..."; \
		sleep 1; \
	done

recreate-db: postgres
	docker exec $$(docker ps -q -f name=postgres) dropdb -U $(DB_USER) --if-exists $(DB_NAME)
	docker exec $$(docker ps -q -f name=postgres) createdb -U $(DB_USER) $(DB_NAME)

migrate: recreate-db
	@for file in $(MIGRATIONS_DIR)/*.up.sql; do \
		echo "Applying $$file..."; \
		docker exec -i $$(docker ps -q -f name=postgres) psql -U $(DB_USER) -d $(DB_NAME) < $$file; \
	done

seed:
	@echo "Seeding database..."
	docker exec -i $$(docker ps -q -f name=postgres) psql -U $(DB_USER) -d $(DB_NAME) < scripts/seed.sql
	@echo "Database seeded successfully!"

restart-app:
	$(DC) restart app

start: migrate seed build up

reset-db: recreate-db migrate seed

.DEFAULT_GOAL := start
