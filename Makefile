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

test:
# シェルコマンドでユニットテスト専用ファイルを順次実行
	docker compose run --rm app bash -c ' \
    array=(`find . -name "*_test.go"`); \
    echo テスト対象ファイル; \
    IFS=$$'\''\n'\''; \
    echo "$${array[*]}"; \
    packages=$$(find . -name "*_test.go" -exec dirname {} \; | sort -u); \
    go test $$packages \
	'
	make down

lint:
	docker compose run --rm app staticcheck ./...

# coverage:
# 	docker compose run --rm app go test ./... -coverprofile=coverage.out
# 	docker compose run --rm app go tool cover -html=coverage.out -o coverage.html