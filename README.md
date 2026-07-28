
# 📊 Pulse: Мультисорсная платформа агрегации и аналитики на Go

> ⚠️ **Статус проекта:** Находится в стадии активной разработки (Реализовано ~85% архитектурного плана).

Микросервисный проект для real-time сбора, нормализации и аналитики событий из разнородных внешних источников (GitHub, Reddit, OpenWeatherMap). <br>
Демонстрирует построение высоконагруженных распределенных систем, продвинутые паттерны конкурентности (Worker Pool, Fan-in, Fan-out) и строгую Чистую Архитектуру (на базе шаблона Evrone).

## 🚀 О проекте

Приложение собирает данные из трех разных API, нормализует их в единый формат `Entity` и пропускает через асинхронный пайплайн обработки (RabbitMQ), данные параллельно записываются в OLAP-хранилище (ClickHouse) для сложной аналитики и в in-memory кэш (Redis) для быстрых счетчиков.

Состоит из пяти независимых микросервисов, общающихся по протоколу **gRPC**:
1. **API Gateway (REST):** HTTP-шлюз на фреймворке Fiber v2. Принимает запросы, проверяет JWT-токены (Middleware) и маршрутизирует их во внутренние gRPC-сервисы.
2. **Collector Service:** Поллит внешние API (GitHub, Reddit, Weather). Реализует паттерн *Fan-in* для объединения данных в единый канал.
3. **Ingestion Service:** gRPC-сервер, принимающий события. Валидирует и публикует их в брокер сообщений RabbitMQ с ручным подтверждением.
4. **Processor Service:** Фоновый демон. Использует *Worker Pool* для чтения из RabbitMQ и паттерн *Fan-out* для параллельной записи батчами (Accumulator) в ClickHouse и Redis.
5. **Analytics Service:** Выполняет тяжелые агрегационные запросы (тренды, топы) в ClickHouse, используя *Cache-Aside* паттерн с Redis.

## 🛠️ Стек технологий

# Backend & Infrastructure
<div> 
<img src="https://img.shields.io/badge/Go_1.22+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go"/> 
<img src="https://img.shields.io/badge/gRPC-244C5A?style=flat&logo=grpc&logoColor=white" alt="gRPC"/>
<img src="https://img.shields.io/badge/Fiber_v2-000000?style=flat&logo=go&logoColor=white" alt="Fiber"/>
<img src="https://img.shields.io/badge/RabbitMQ-FF6600?style=flat&logo=rabbitmq&logoColor=white" alt="RabbitMQ"/>
<img src="https://img.shields.io/badge/ClickHouse-FFCC01?style=flat&logo=clickhouse&logoColor=black" alt="ClickHouse"/>
<img src="https://img.shields.io/badge/PostgreSQL-336791?style=flat&logo=postgresql&logoColor=white" alt="PostgreSQL"/> 
<img src="https://img.shields.io/badge/Redis-DC382D?style=flat&logo=redis&logoColor=white" alt="Redis"/>
<img src="https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white" alt="Docker"/> 
</div>

## ⚙️ Как запустить локально

### Необходимые компоненты
* [Go](https://golang.org/dl/) (версия 1.22+)
* [Docker и Docker Compose](https://www.docker.com/)
* `make` (утилита сборки)

### 1. Клонируйте репозиторий

```sh
git clone https://github.com/Maxim-Belyi/pulse.git
cd Pulse
```


> **Важно:** Скрипты миграций автоматически создадут необходимые таблицы в PostgreSQL и ClickHouse при старте.

## 🌐 Архитектура и паттерны

```text
[Внешние API (GitHub, Reddit...)] ──► [Collector Service] (Поллинг, Fan-in)
                                                │
                                                ▼ (gRPC)
[Клиент (HTTP/REST)] ──► [API Gateway] ──► [Ingestion Service] ──► [RabbitMQ]
                             │                                         │
                             ▼ (gRPC)                                  ▼ (Worker Pool)
                    [Analytics Service] ◄──(gRPC)────────────── [Processor Service] 
                             │                                (Fan-out, Батчинг)
                             ▼ (Cache-Aside)                           │
                       [ClickHouse & Redis] ◄──────────────────────────┘
```

* **Чистая Архитектура (Evrone Template):** Строгое разделение на слои `entity`, `usecase`, `repo` и `controller`. Зависимости направлены внутрь, инфраструктурный слой ничего не знает о бизнес-логике и общается через интерфейсы.
* **Concurrency:** Работа с горутинами, контекстами (`context.Context`), каналами и `sync.WaitGroup`, реализован механизм Graceful Shutdown для безопасного завершения работы с ожиданием сброса батчей.
* **Batching:** Запись в ClickHouse происходит пачками (по размеру или таймеру) для минимизации нагрузки на диск.

## 🗺️ Roadmap (В планах на ближайшие релизы)

- [x] Интеграция всех источников данных
- [x] Асинхронный конвейер обработки и батчинг
- [x] gRPC контроллеры и аналитический движок
- [ ] Финализация API Gateway и Swagger-документации 
- [ ] Написание изолированных Unit-тестов (testify + mockgen)
- [ ] Интеграционные тесты с Testcontainers
- [ ] Добавление распределенного трейсинга OpenTelemetry