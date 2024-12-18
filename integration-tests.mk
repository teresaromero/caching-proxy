BINARY_NAME=integration/caching-proxy

.PHONY: test-integration clean-integration
test-integration:
	@echo "Running tests..."
	@go build -o $(BINARY_NAME) . \
	&& (trap '$(MAKE) clean-integration' EXIT; go test -v -tags=integration ./integration/...)
	
clean-integration:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)