include app.env


sqlc:
	sqlc generate
migrateup:
	goose -dir db/migration $(DB_DRIVER) $(DB_URL) up
migratedown:
	goose -dir db/migration $(DB_DRIVER) $(DB_URL) down
test:
	go test -v -cover -short ./...
server:
	go build -o server main.go RDS_DB_URL=$(RDS_DB_URL) RDS_DB_USER=$(RDS_DB_USER) RDS_DB_PASSWORD=$(RDS_DB_PASSWORD)


.PHONY: sqlc migrateup migratedown test server
