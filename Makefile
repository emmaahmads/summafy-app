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
	go build -o server main.go
	./server


.PHONY: sqlc migrateup migratedown test server
