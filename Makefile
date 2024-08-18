include .env

sqlccompile:
	sqlc compile
sqlcgen:
	sqlc generate
migrateup:
	goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up
migratedown:
	goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down
test:
	go test -v -cover -short ./...

server:
	go run main.go

.PHONY: sqlccompile sqlcgen migrateup migratedown test server
