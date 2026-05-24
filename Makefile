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
	@echo "  migrate-create NAME=<migration_name> - Create a new migration file"
	@echo "  migrate-up   - Apply all up migrations"
	@echo "  migrate-down - Apply all down migrations"
	@echo "  migrate-force VERSION=<version> - Force set the migration version"
	@echo "  help    - Show this help message"
