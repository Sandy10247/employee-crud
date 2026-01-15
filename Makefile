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