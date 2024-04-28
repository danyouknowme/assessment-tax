POSTGRES_URL=postgres://postgres:postgres@localhost:5432/ktaxes?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose down

mock:
	mockgen -package=mockdb -source="db/store.go" -destination="db/mock/store.go" Store

run:
	go run main.go

test:
	go test -v -cover ./...

it-test:
	go test -v -tags=integration ./...

.PHONY: migrateup mock run test it-test