include .env
export

DIR := ${CURDIR}
MIGRATE := docker run -v $(DIR)/migrations:/migrations \
		   --network host migrate/migrate -path=/migrations/ -database "$(POSTGRES_CONN)?sslmode=disable"

.PHONY: migrate migrate-down testdata
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

testdata: ## populate the database with test data
	@echo "Populating test data..."
	@docker exec -i postgres psql "$(POSTGRES_CONN)?sslmode=disable" < ./testdata/testdata.sql
