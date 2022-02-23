#!make

.PHONY: run
run:
	docker-compose run --rm app go run . $(FROM) $(TO) $(DIRNAME)

.PHONY: test
test:
	docker-compose run --rm app go test ./converter

.PHONY: clean
clean:
	rm -rf ./result/*
