include .env

sqlc:
	sqlc generate
migrateup:
	goose -dir db/migration $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up
migratedown:
	goose -dir db/migration $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down
test:
	go test -v -cover -short ./...

server:
	go run main.go

.PHONY: sqlc migrateup migratedown test server
