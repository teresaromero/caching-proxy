.PHONY: origin-run origin-stop origin-clean

origin-run:
	@echo "Starting Origin container..."
	@docker run -d \
		--name origin \
		-p 8080:5678 \
		-text="Hello, world!" \
		hashicorp/http-echo

origin-stop:
	@echo "Stopping Origin container..."
	@docker stop origin || true
	@docker rm origin || true

origin-clean:
	@echo "Removing container and volume..."
	@docker stop origin || true
	@docker rm origin || true