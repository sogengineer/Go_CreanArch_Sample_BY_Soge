.PHONY: all

ifeq ($(OS),Windows_NT)
    CP = copy
	COPY = cmd /c "move /y src\\coverage.* ."
	OPEN = cmd /c "start coverage.html"
else
    CP = cp
	COPY = mv src/coverage.* ./ -f
	UNAME := $(shell uname)
    ifeq ($(UNAME),Linux)
		OPEN = cmd.exe /c start coverage.html
    endif
    ifeq ($(UNAME),Darwin)
        OPEN = open coverage.html
    endif
endif

all:
	make cleanup
	make build

build:
	$(CP) env.example .env
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
	docker compose down

lint:
	docker compose run --rm app staticcheck ./...

coverage:
	docker compose run --rm app bash -c ' \
		array=(`find . -name "*_test.go"`); \
		echo テスト対象ファイル; \
		IFS=$$'\''\n'\''; \
		echo "$${array[*]}"; \
		packages=$$(find . -name "*_test.go" -exec dirname {} \; | sort -u); \
        go test $${packages[@]} -coverprofile=coverage.out; \
        go tool cover -html=coverage.out -o coverage.html \
        '
	$(COPY)
	$(OPEN)