default: lint

lint:
	golangci-lint run --config ci/golangci-lint/golangci.yml --timeout=5m
