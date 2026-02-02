DB_DSN=postgres://postgres:12345@auth-db:5432/authorization?sslmode=disable

MIGRATE := migrate -path ./migrations -database $(DB_DSN)

migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

lint:
	golangci-lint run