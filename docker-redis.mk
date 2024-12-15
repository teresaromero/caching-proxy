.PHONY: redis-run redis-stop redis-shell redis-clean

redis-run:
	@echo "Starting Redis container..."
	@docker run -d \
		--name my-redis \
		-p 127.0.0.1:6379:6379 \
		-v redis-data:/data \
		redis:latest \
		redis-server --requirepass "password"

redis-stop:
	@echo "Stopping Redis container..."
	@docker stop my-redis || true
	@docker rm my-redis || true

redis-logs:
	@echo "Tailing logs from Redis container..."
	@docker logs -f my-redis

redis-shell:
	@echo "Attaching redis-cli to the running Redis container..."
	@docker exec -it my-redis redis-cli

redis-clean:
	@echo "Removing container and volume..."
	@docker stop my-redis || true
	@docker rm my-redis || true
	@docker volume rm redis-data || true