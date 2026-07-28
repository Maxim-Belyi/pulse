COMPOSE_FILE := deploy/docker-compose.yml
PROTO_SRC    := docs/proto/v1
PROTO_OUT    := .
SERVICES     := collector ingestion processor analytics api-gateway
TEST_FLAGS   := -v -race -count=1

.PHONY: help \
        compose-up compose-down compose-restart compose-logs compose-ps compose-clean \
        proto \
        build build-all \
        test test-cover test-integration \
        lint lint-fix \
        mock \
        migrate-up migrate-down migrate-create \
        tidy vet clean


help: ## Показать справку по командам
	@echo ""
	@echo "Pulse — команды разработки"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) \
		| awk -F':.*##' '{printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""

compose-up: ## Запустить весь стек (PostgreSQL, ClickHouse, Redis, RabbitMQ)
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "  Стек запущен."
	@echo "  RabbitMQ UI : http://localhost:15672  (guest/guest)"
	@echo "  ClickHouse  : http://localhost:8123"
	@echo "  PostgreSQL  : localhost:5432"
	@echo "  Redis       : localhost:6379"

compose-down: ## Остановить контейнеры (данные сохраняются)
	docker compose -f $(COMPOSE_FILE) down

compose-restart: ## Перезапустить весь стек
	docker compose -f $(COMPOSE_FILE) restart

compose-logs: ## Логи всех контейнеров в реальном времени
	docker compose -f $(COMPOSE_FILE) logs -f

compose-ps: ## Статус контейнеров
	docker compose -f $(COMPOSE_FILE) ps

compose-clean: ## Остановить контейнеры и удалить тома (удалит данные!)
	docker compose -f $(COMPOSE_FILE) down -v

proto: ## Сгенерировать Go-код из .proto файлов
	protoc \
		--go_out=$(PROTO_OUT) \
		--go_opt=module=pulse \
		--go-grpc_out=$(PROTO_OUT) \
		--go-grpc_opt=module=pulse \
		$(PROTO_SRC)/*.proto
	@echo "  proto сгенерированы -> api/pb/"

build: ## Собрать один сервис: make build SERVICE=collector
	@if [ -z "$(SERVICE)" ]; then echo "Укажи: make build SERVICE=collector"; exit 1; fi
	go build -o bin/$(SERVICE) ./services/$(SERVICE)/cmd/app/
	@echo "  Собран: bin/$(SERVICE)"

build-all: ## Собрать все 5 сервисов в ./bin/
	@for svc in $(SERVICES); do \
		echo "-> Сборка $$svc..."; \
		go build -o bin/$$svc ./services/$$svc/cmd/app/ || exit 1; \
	done
	@echo "  Все сервисы собраны в ./bin/"

test: ## Запустить unit-тесты с -race
	go test $(TEST_FLAGS) ./internal/...

test-cover: ## Тесты + HTML-отчёт о покрытии (coverage.html)
	go test $(TEST_FLAGS) -coverprofile=coverage.txt -covermode=atomic ./internal/...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "  Отчёт: coverage.html"

test-integration: ## Интеграционные тесты (требует Docker)
	go test $(TEST_FLAGS) -tags=integration ./...

lint: ## Запустить golangci-lint
	golangci-lint run ./...

lint-fix: ## Запустить golangci-lint с автоисправлением
	golangci-lint run --fix ./...

mock: ## Сгенерировать моки из usecase/contracts.go
	mockgen -source=internal/usecase/contracts.go \
	        -destination=internal/usecase/mocks/mocks.go \
	        -package=mocks
	@echo "  Моки сгенерированы -> internal/usecase/mocks/"

migrate-up: ## Применить все миграции PostgreSQL
	migrate -path migrations/postgres -database "$(POSTGRES_DSN)" up

migrate-down: ## Откатить последнюю миграцию PostgreSQL
	migrate -path migrations/postgres -database "$(POSTGRES_DSN)" down 1

migrate-create: ## Создать миграцию: make migrate-create NAME=create_events
	@if [ -z "$(NAME)" ]; then echo "Укажи: make migrate-create NAME=create_events"; exit 1; fi
	migrate create -ext sql -dir migrations/postgres -seq $(NAME)
	@echo "  Создана миграция в migrations/postgres/"

tidy: ## Привести go.mod и go.sum в порядок
	go mod tidy

vet: ## Запустить go vet
	go vet ./...

clean: ## Удалить артефакты сборки (bin/, coverage.*)
	rm -rf bin/ coverage.txt coverage.html