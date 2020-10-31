DB_NAME=db.sqlite

migrate_up:
	goose -dir migrations sqlite3 $(DB_NAME) up

migrate_down:
	goose -dir migrations sqlite3 $(DB_NAME) down

sqlc:
	sqlc generate

go-test:
	go test ./test/... -bench=.
