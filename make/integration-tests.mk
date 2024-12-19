BINARY_NAME=integration/caching-proxy

.PHONY: integration-tests clean-integration
integration-tests:
	@echo "Running integration tests..."
	@go build -o $(BINARY_NAME) . \
	&& (trap '$(MAKE) clean-integration' EXIT; go test -v -tags=integration ./integration/...)
	
clean-integration:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)