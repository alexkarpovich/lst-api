$PWD=$(shell pwd)

.PHONY: create_migration apply_migration

create_migration:
	migrate create -ext sql -dir $(PWD)/src/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

apply_migration:
	migrate -path src/migrations -database "postgres://postgres:postgres@0.0.0.0:5432/dev?sslmode=disable" $(filter-out $@,$(MAKECMDGOALS))
