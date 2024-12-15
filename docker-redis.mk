.PHONY: redis-run redis-stop redis-shell redis-clean

redis-run:
	@echo "Starting Redis container..."
	@docker run -d \
		--name my-redis \
		-p 6379:6379 \
		-v redis-data:/data \
		redis:latest

redis-stop:
	@echo "Stopping Redis container..."
	@docker stop $(CONTAINER_NAME) || true
	@docker rm $(CONTAINER_NAME) || true

redis-logs:
	@echo "Tailing logs from Redis container..."
	@docker logs -f $(CONTAINER_NAME)

redis-shell:
	@echo "Attaching redis-cli to the running Redis container..."
	@docker exec -it $(CONTAINER_NAME) redis-cli

redis-clean:
	@echo "Removing container and volume..."
	@docker stop $(CONTAINER_NAME) || true
	@docker rm $(CONTAINER_NAME) || true
	@docker volume rm $(VOLUME_NAME) || true