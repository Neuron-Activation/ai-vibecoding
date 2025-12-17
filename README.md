# Go Application — Notes Service (обновлённый README)

Сервис для работы с заметками на языке Go.

---

## Кратко

Проект — минимальный REST-сервис для CRUD-операций с заметками (notes). Использует:

* Go
* Gorilla Mux (роутинг)
* GORM (ORM)
* PostgreSQL (в production) или SQLite (для разработки и тестов)

---

## Что изменено (важно)

* Рефакторинг `db`-пакета: убрана автоматическая инициализация в `init()`. Добавлены функции:

  * `InitDB(dialect, conn string) error` — явная инициализация
  * `InitDBFromEnv() error` — инициализация по переменным окружения
  * `CloseDB() error` — закрытие соединения

  Это позволяет запускать приложение с in-memory SQLite в тестах/локальной разработке и упрощает покрытие тестами контроллеров.

* `main.go` обновлён: теперь вызывает `db.InitDBFromEnv()` и корректно закрывает соединение с БД.

* Добавлен README с инструкциями запуска, Docker Compose примером и рекомендациями по тестированию.

---

## Структура проекта

```
go-app/
  controllers/    # HTTP handlers
  db/             # DB init & миграции
  models/         # модели GORM
  utils/          # хелперы и обработка ошибок
  main.go
  go.mod
  README.md
```

---

## Переменные окружения

* Для Postgres (по умолчанию):

  * `db_user` (пример: `postgres`)
  * `db_pass` (пример: `postgres`)
  * `db_name` (пример: `postgres`)
  * `db_host` (пример: `localhost` или `db` в docker-compose)
  * `db_port` (пример: `5432`)

* Опционально можно указать `DB_DIALECT`:

  * `postgres` (по умолчанию)
  * `sqlite3` — для разработки/тестов

* Для sqlite3 укажите `DB_CONN` (например `:memory:` или `file:test.db?cache=shared`).

---

## Локальный запуск (Windows / PowerShell)

1. Установите Go (рекомендуется 1.18+). Убедитесь, что `go` в PATH.
2. Скопируйте `.env-example` в `.env` или задайте переменные окружения через PowerShell:

```powershell
$env:DB_DIALECT = "postgres"
$env:db_user = "postgres"
$env:db_pass = "postgres"
$env:db_name = "postgres"
$env:db_host = "localhost"
$env:db_port = "5432"
```

3. Соберите и запустите:

```powershell
go mod tidy
go build -o main .
.\\main.exe
```

Сервис будет доступен по адресу: `http://localhost:8080`.

### Быстрый запуск без Postgres — SQLite in-memory (для разработки)

В PowerShell перед запуском:

```powershell
$env:DB_DIALECT = "sqlite3"
$env:DB_CONN = ":memory:"
go run main.go
```

Это удобно для быстрого старта и для тестов контроллеров.

---

## Запуск через Docker (Postgres + приложение)

### docker-compose.yml (пример)

```yaml
version: '3.8'
services:
  db:
    image: postgres:13
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: postgres
    ports:
      - "5432:5432"
    volumes:
      - db-data:/var/lib/postgresql/data

  app:
    build: .
    environment:
      DB_DIALECT: postgres
      db_user: postgres
      db_pass: postgres
      db_name: postgres
      db_host: db
      db_port: "5432"
    ports:
      - "8080:8080"
    depends_on:
      - db

volumes:
  db-data:
```

1. `docker-compose up -d`
2. `docker-compose logs -f app` (или `docker-compose exec app /bin/sh` для интерактивного доступа)

> Примечание: Docker image приложения можно собрать через `docker build -t go-app .`.

---

## API

* `GET /notes` — получить список. Опциональный параметр `?query=<text>` ищет по title/content.
* `POST /notes` — создать заметку. JSON: `{ "title": "..", "content": ".." }`.
* `GET /notes/{id}` — получить по id.
* `PUT /notes/{id}` — обновить.
* `DELETE /notes/{id}` — удалить.

Пример curl:

```bash
curl -X POST http://localhost:8080/notes -d '{"title":"t","content":"c"}' -H "Content-Type: application/json"
```

---

## Тестирование

### Unit tests (локально)

Тесты для `utils` не зависят от БД и запускаются так:

```bash
go test ./utils -v
```

Общий запуск всех тестов:

```bash
go test ./... -v
```

### Тесты контроллеров / интеграционные тесты

Благодаря `db.InitDB(dialect, conn)` в тестах вы можете инициализировать in-memory sqlite и прогонять контроллеры через `httptest`.

Пример в `TestMain`:

```go
func TestMain(m *testing.M) {
    _ = db.InitDB("sqlite3", ":memory:")
    code := m.Run()
    _ = db.CloseDB()
    os.Exit(code)
}
```

Это позволит делать реальные CRUD-операции в тестах без поднятия Postgres.

---

## Рекомендации по улучшению

* Вынести слой доступа к данным (repository) за интерфейсами — упростит мокинг и unit-тестирование.
* Использовать управляемые миграции (golang-migrate) вместо `AutoMigrate` для production.
* Добавить CI (GitHub Actions): сборка, `go test ./...`, линтер, сборка Docker image.
* Добавить аналитические эндпоинты и метрики (Prometheus) — полезно для нагрузочных тестов.

---

Если нужно — могу:

* Сгенерировать patch/PR с изменениями в `db` и `main.go`;
* Добавить файл `docker-compose.yml` в репозиторий;
* Сгенерировать тесты контроллеров и показать пример `TestMain` и пару тестов CRUD.

---

Спасибо — скажите, что делаем дальше?
