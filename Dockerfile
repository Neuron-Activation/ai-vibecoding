FROM golang:1.20-bullseye

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    sqlite3 \
    libsqlite3-dev \
 && rm -rf /var/lib/apt/lists/*

# Рабочая директория
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем весь проект
COPY . .

# Собираем бинарник в отдельную папку, чтобы volume не затирал его
RUN go build -o /go/bin/main main.go

# Экспонируем порт
EXPOSE 8080

# Команда запуска
CMD ["/go/bin/main"]
