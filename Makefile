.PHONY: run

run:
	@echo "Running server in debug mode..."
	@DEBUG=true go run cmd/gokku/main.go
