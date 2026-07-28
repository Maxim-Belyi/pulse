COMPOSE_FILE := deploy/docker-compose.yml

PROTO_SRC    := docs/proto/v1
PROTO_OUT    := .

SERVICES := collector ingestion processor analytics api-gateway

TEST_FLAGS := -v -race -count=1

.PHONY: help \
        compose-up compose-down compose-restart compose-logs compose-ps \
        proto \
        build build-all \
        test test-cover lint \
        mock \
        migrate-up migrate-down \
        tidy \
        clean

help: 
	@echo ""
	@echo "Pulse — команды разработки"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""


compose-up: 
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "    Стек запущен."
	@echo "    RabbitMQ UI : http://localhost:15672  (guest/guest)"
	@echo "    ClickHouse  : http://localhost:8123"
	@echo "    PostgreSQL  : localhost:5432"
	@echo "    Redis       : localhost:6379"

compose-down: 
	docker compose -f $(COMPOSE_FILE) down

compose-restart:
	docker compose -f $(COMPOSE_FILE) restart

compose-logs: 
	docker compose -f $(COMPOSE_FILE) logs -f

compose-ps:
	docker compose -f $(COMPOSE_FILE) ps

compose-clean: 
	docker compose -f $(COMPOSE_FILE) down -v

proto: 
	protoc \
		--go_out=$(PROTO_OUT) \
		--go_opt=module=pulse \
		--go-grpc_out=$(PROTO_OUT) \
		--go-grpc_opt=module=pulse \
		$(PROTO_SRC)/*.proto
	@echo "proto сгенерированы → api/pb/"

build: 
	@if [ -z "$(SERVICE)" ]; then echo "Укажи: make build SERVICE=collector"; exit 1; fi
	go build -o bin/$(SERVICE) ./services/$(SERVICE)/cmd/app/
	@echo "Собран: bin/$(SERVICE)"

build-all:
	@for svc in $(SERVICES); do \
		echo "→ Сборка $$svc..."; \
		go build -o bin/$$svc ./services/$$svc/cmd/app/ || exit 1; \
	done
	@echo "Все сервисы собраны в ./bin/"


test:
	go test $(TEST_FLAGS) ./internal/...
	@echo "Тесты пройдены"

test-cover:
	go test $(TEST_FLAGS) -coverprofile=coverage.txt -covermode=atomic ./internal/...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Отчёт: coverage.html"

test-integration: 
	go test $(TEST_FLAGS) -tags=integration ./...
	@echo "Интеграционные тесты пройдены"


lint: 
	golangci-lint run ./...

lint-fix: 
	golangci-lint run --fix ./...

mock:
	mockgen -source=internal/usecase/contracts.go \
	        -destination=internal/usecase/mocks/mocks.go \
	        -package=mocks
	@echo "Моки сгенерированы → internal/usecase/mocks/"


migrate-up: 
	migrate -path migrations/postgres -database "$(POSTGRES_DSN)" up

migrate-down:
	migrate -path migrations/postgres -database "$(POSTGRES_DSN)" down 1

migrate-create: 
	@if [ -z "$(NAME)" ]; then echo "Укажи: make migrate-create NAME=create_events"; exit 1; fi
	migrate create -ext sql -dir migrations/postgres -seq $(NAME)
	@echo "Создана миграция: migrations/postgres/"


tidy: 
	go mod tidy
	@echo "go mod tidy выполнен"

vet: 
	go vet ./...

clean: 
	rm -rf bin/ coverage.txt coverage.html
	@echo "Очищено"