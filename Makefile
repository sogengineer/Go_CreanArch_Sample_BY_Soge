.PHONY: all

all:
	docker compose down
	docker system prune --volumes -f
	docker compose build
	docker compose run --rm app go mod tidy
	docker compose up -d --force-recreate

build:
	docker compose build
	docker compose run --rm app go mod tidy
	docker compose up -d --force-recreate

cleanup:
	docker compose down
	docker system prune --volumes -f

up:
	docker compose build
	docker compose run --rm app go mod tidy
	docker compose up -d

down:
	docker compose down
