build:
	@go fmt ./...
	@go build -o bin/gravemind

run: build
	@./bin/gravemind

up:
	@podman compose up -d

down:
	@podman compose down
