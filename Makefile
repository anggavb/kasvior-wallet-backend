include ./.env

MIGRATION_PATH=db/migrations
SEEDER_PATH=db/seeds
SEED_TABLE=schema_seeds

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

seeder-create:
	@touch $(SEEDER_PATH)/$(NAME)_seeder.sql

seed:
	@echo "Start seeding..."
	@for f in ${SEEDER_PATH}/*.sql; do \
		echo "Applying seeder: $$f"; \
		psql $(DB_URL) -f "$$f"; \
	done
	@echo "Seeding completed."

fresh:
	@make migrate-down
	@make migrate-up
	@make seed

help:
	@echo "Available commands:"
	@echo "  build                                  - Build the application"
	@echo "  run                                    - Build and run the application"
	@echo "  fresh                                  - Reset the database and reapply all migrations and seeders"
	@echo "  migrate-create NAME=<migration_name>   - Create a new migration file"
	@echo "  migrate-up                             - Apply all up migrations"
	@echo "  migrate-down                           - Apply all down migrations"
	@echo "  migrate-force VERSION=<version>        - Force set the migration version"
	@echo "  seeder-create NAME=<seeder_name>       - Create a new seeder file"
	@echo "  seed                                   - Apply all seeders"
	@echo "  help                                   - Show this help message"
