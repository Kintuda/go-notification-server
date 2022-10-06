build:
	go build
dev:
	air
create-new-migration:
	go run main.go migration create $(arg1)
migration-up:
	go run main.go migration up