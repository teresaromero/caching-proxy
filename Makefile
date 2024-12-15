include docker-redis.mk
include docker-origin.mk

.PHONY: run environment environment-stop environment-clean
run:
	@echo "Running the Go application..."
	go run main.go --origin=http://localhost:8080 --port=8081

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