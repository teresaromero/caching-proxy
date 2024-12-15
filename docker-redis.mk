REDIS_PASSWORD=password

.PHONY: redis-run redis-stop redis-shell redis-clean

redis-run:
	@echo "Starting Redis container..."
	@docker run -d \
		--name redis \
		-p 127.0.0.1:6379:6379 \
		-v redis-data:/data \
		redis:latest \
		redis-server --requirepass $(REDIS_PASSWORD)

redis-stop:
	@echo "Stopping Redis container..."
	@docker stop redis || true
	@docker rm redis || true

redis-shell:
	@echo "Attaching redis-cli to the running Redis container..."
	@docker exec -it redis redis-cli -a $(REDIS_PASSWORD)

redis-clean:
	@echo "Removing container and volume..."
	@docker stop redis || true
	@docker rm redis || true
	@docker volume rm redis-data || true