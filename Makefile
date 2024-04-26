POSTGRES_URL=postgres://postgres:postgres@localhost:5432/ktaxes?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose up

mock:
	mockgen -package=mockdb -source="db/store.go" -destination="db/mock/store.go" Store

run:
	go run main.go

test:
	go test -v -cover ./...

.PHONY: migrateup mock run test