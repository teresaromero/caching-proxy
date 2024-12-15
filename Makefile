PORT:=8081
ORIGIN_HOST=http://localhost
ORIGIN_PORT=8080

REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=password

include docker-redis.mk
include docker-origin.mk

.PHONY: run environment environment-stop environment-clean
run:
	@echo "Running the Go application..."
	go run main.go --origin=$(ORIGIN_HOST):$(ORIGIN_PORT) --port=$(PORT)

environment:
	@echo "Running redis and origin containers..."
	@make redis-run
	@make origin-run

environment-stop:
	@echo "Stopping redis and origin containers..."
	@make redis-stop
	@make origin-stop

environment-clean:
	@echo "Removing redis and origin containers..."
	@make redis-clean
	@make origin-clean