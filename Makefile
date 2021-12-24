$PWD=$(shell pwd)

.PHONY: m

m:
ifdef new
	migrate create -ext sql -dir $(PWD)/src/migrations -seq $(new)
else 
	COMMAND = 'migrate -path src/migrations -database "postgres://postgres:postgres@0.0.0.0:5432/dev?sslmode=disable" '
	ifdef up
		COMMAND += 'up'
	else ifdef down
		COMMAND += 'down'
	else ifdef force
		COMMAND += 'force'
	endif

	$(COMMAND)
endif
