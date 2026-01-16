# Define ANSI color codes as variables
COL_RED=\033[0;31m
COL_GREEN=\033[0;32m

run: 
	@echo "Running Server....."
	@go run cmd/*.go

instal_sqlc :
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

install_goose :
	@go install github.com/pressly/goose/v3/cmd/goose@latest

goose_status:
	cd sql/schema && go build -o goose-custom *.go && ./goose-custom status && rm goose-custom

goose_up:
	cd sql/schema && go build -o goose-custom *.go && ./goose-custom up && rm goose-custom

goose_down:
	cd sql/schema && go build -o goose-custom *.go && ./goose-custom down && rm goose-custom

sqlc:
	sqlc generate

clean_shit:
	@echo "Removing all branches Except `master`"
	@git checkout master && git branch | grep -v "master" | xargs git branch -D 
	@echo "$(COL_GREEN)Pull `master`$(COL_GREEN)"
	@git pull

docker_up:
	@echo "$(COL_GREEN)Running Docker Compose ⬆️ ✅$(COL_GREEN)"
	@docker compose up -d --build

docker_down:
	@echo "$(COL_RED)Running Docker Compose ⬇️ 🔴 $(COL_RED)"
	@docker compose down
