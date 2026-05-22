include ./.env

MIGRATION_PATH=db/migrations

build:
	@go run build ./cmd/main.go

run:
	@go run ./cmd/main.go

migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	@migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down:
	@migrate -database $(DB_URL) -path $(MIGRATION_PATH) down

migrate-force:
	@migrate -database $(DB_URL) -path $(MIGRATION_PATH) force $(VERSION)

help:
	@echo "Available commands:"
	@echo "  build   - Build the application"
	@echo "  run     - Build and run the application"
