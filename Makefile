install:
	go mod download
build:
	go build -o bin/main main.go
start:
	go run main.go server
migration_up:
	migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRESQL_URL)" -verbose up
migration_down:
	migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRESQL_URL)" -verbose down
migration_force:
	migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRESQL_URL)" -verbose force $(VERSION)
create_migration:
	migrate create -ext sql -dir $