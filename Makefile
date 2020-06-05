.PHONY: build-go

build-go: ## Build all Go binaries.
	@echo "build go files"
	CGO_ENABLED=0 go build -mod vendor -o ./bin/

build-docker: ## Build Docker images.
	@echo "build docker"
	go mod tidy
	go mod vendor
	docker build -t backend-demo:latest ./