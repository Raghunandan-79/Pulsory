# Makefile
.PHONY: help build up down restart logs clean ps

help:
	@echo "Available commands:"
	@echo "  make build    - Build all Docker images"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make restart  - Restart all services"
	@echo "  make logs     - Show logs from all services"
	@echo "  make logs-server - Show server logs"
	@echo "  make logs-pusher - Show pusher logs"
	@echo "  make logs-worker - Show worker logs"
	@echo "  make ps       - Show running containers"
	@echo "  make clean    - Remove containers, volumes, and images"

build:
	docker-compose build

up:
	docker-compose up -d
	@echo "All services started!"

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

logs-server:
	docker-compose logs -f server

logs-pusher:
	docker-compose logs -f pusher

logs-worker:
	docker-compose logs -f worker

ps:
	docker-compose ps

clean:
	docker-compose down -v
	docker-compose rm -f
	docker system prune -f

# Development commands
dev-server:
	go run cmd/server/main.go

dev-pusher:
	go run cmd/pusher/main.go

dev-worker:
	REGION_ID=region-1 WORKER_ID=worker-1 go run cmd/worker/main.go

dev-all:
	go run cmd/server/main.go & \
	go run cmd/pusher/main.go & \
	REGION_ID=region-1 WORKER_ID=worker-1 go run cmd/worker/main.go & \
	wait

# Production build
build-local:
	go build -o bin/server cmd/server/main.go
	go build -o bin/pusher cmd/pusher/main.go
	go build -o bin/worker cmd/worker/main.go

run-local:
	./bin/server & \
	./bin/pusher & \
	./bin/worker & \
	wait