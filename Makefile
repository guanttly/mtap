# MTAP Makefile
# 核心目的：项目构建、测试、部署的快捷命令

.PHONY: build-api build-sync build-scheduler build-ws build-notify build-analyzer build-all \
        migrate-up migrate-down lint test docker-up docker-down

build-api:
go build -o bin/api-server ./cmd/api-server

build-sync:
go build -o bin/sync-worker ./cmd/sync-worker

build-scheduler:
go build -o bin/scheduler ./cmd/scheduler

build-ws:
go build -o bin/ws-server ./cmd/ws-server

build-notify:
go build -o bin/notify-worker ./cmd/notify-worker

build-analyzer:
go build -o bin/analyzer-worker ./cmd/analyzer-worker

build-all: build-api build-sync build-scheduler build-ws build-notify build-analyzer

migrate-up:
go run ./cmd/migrate up

migrate-down:
go run ./cmd/migrate down

lint:
golangci-lint run ./...

test:
go test ./... -cover

docker-up:
docker-compose -f deployments/docker-compose.yaml up -d

docker-down:
docker-compose -f deployments/docker-compose.yaml down
