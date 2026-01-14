run: 
	@echo "Running Server....."
	@go run cmd/*.go

instal_sqlc :
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

install_goose :
	@go install github.com/pressly/goose/v3/cmd/goose@latest

goose_up:
	cd sql/schema && goose postgres postgres://phani:postgres@localhost:5432/test1 up

goose_down:
	cd sql/schema && goose postgres postgres://phani:postgres@localhost:5432/test1 down

sqlc:
	sqlc generate