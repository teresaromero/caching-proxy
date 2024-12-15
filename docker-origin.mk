.PHONY: origin-run origin-stop origin-shell origin-clean

origin-run:
	@echo "Starting Origin container..."
	@docker run -d \
		--name my-origin \
		-p 8080:5678 \
		-text="Hello, world!" \
		hashicorp/http-echo

origin-stop:
	@echo "Stopping Origin container..."
	@docker stop $(CONTAINER_NAME) || true
	@docker rm $(CONTAINER_NAME) || true

origin-logs:
	@echo "Tailing logs from Origin container..."
	@docker logs -f $(CONTAINER_NAME)

origin-shell:
	@echo "Attaching origin-cli to the running Origin container..."
	@docker exec -it $(CONTAINER_NAME) origin-cli

origin-clean:
	@echo "Removing container and volume..."
	@docker stop $(CONTAINER_NAME) || true
	@docker rm $(CONTAINER_NAME) || true
	@docker volume rm $(VOLUME_NAME) || true