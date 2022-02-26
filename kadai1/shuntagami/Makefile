#!make

.PHONY: run
run:
	docker-compose run --rm app go run . $(FROM) $(TO) $(DIRNAME)

.PHONY: test
test:
	docker-compose run --rm app go test ./converter

.PHONY: cover
cover:
	docker-compose run --rm app go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html

.PHONY: clean
clean:
	rm -rf ./result/*
