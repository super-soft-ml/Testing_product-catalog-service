# Product Catalog Service

.PHONY: deps migrate test run proto docker-up docker-down

# Spanner emulator
SPANNER_PROJECT  ?= test-project
SPANNER_INSTANCE ?= test-instance
SPANNER_DATABASE ?= product-catalog
SPANNER_EMULATOR_HOST ?= localhost:9010

export SPANNER_EMULATOR_HOST

deps:
	go mod download
	go mod tidy

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/product/v1/product_service.proto

migrate:
	@echo "Ensure Spanner emulator is running: docker-compose up -d"
	@which gcloud >/dev/null 2>&1 || (echo "Install Google Cloud SDK (gcloud) for migrations"; exit 1)
	gcloud config set project $(SPANNER_PROJECT)
	gcloud emulators spanner start --host-port=localhost:9010 &
	@sleep 3
	@echo "Create instance/database if needed, then run DDL from migrations/001_initial_schema.sql"
	@echo "On Windows, run migrations manually - see README."

test:
	go test ./...

test-e2e:
	$(eval export SPANNER_EMULATOR_HOST=localhost:9010)
	go test -v -tags=e2e ./tests/e2e/...

run:
	go run ./cmd/server

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
