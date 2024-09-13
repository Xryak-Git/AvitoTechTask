.PHONY: migrate lint test


migrate:
	migrate -path migrations -database '$(POSTGRES_CONN)?sslmode=disable' up

migrate-up: ### migration up
	migrate -path /migrations -database '$(POSTGRES_CONN)?sslmode=disable' up
.PHONY: migrate-up