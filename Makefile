DB_NAME=db.sqlite
SWAGGER_PATH=controllers/api/swagger.yml
UI_PORT=8081

migrate_up:
	goose -dir migrations sqlite3 $(DB_NAME) up

migrate_down:
	goose -dir migrations sqlite3 $(DB_NAME) down

sqlc:
	sqlc generate

go-test:
	go test ./test/... -bench=.

swagger:
	swagger generate spec -w ./controllers -o $(SWAGGER_PATH)

swagger-ui:
	swagger serve -p $(UI_PORT) --flavor=swagger $(SWAGGER_PATH)
