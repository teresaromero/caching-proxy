include docker-redis.mk
include docker-origin.mk

.PHONY: run
run:
	@echo "Running the Go application..."
	go run main.go --origin=http://localhost:8080 --port=8081