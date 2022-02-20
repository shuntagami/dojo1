#!make

.PHONY: run
run:
	docker-compose run --rm app go run . $(FROM) $(TO) $(DIRNAME)

.PHONY: clean
clean:
	rm -rf ./converted/*
