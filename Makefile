POSTGRES_URL=postgres://postgres:postgres@localhost:5432/ktaxes?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose up

.PHONY: migrateup