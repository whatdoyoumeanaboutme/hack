FROM golang:1.20-alpine

# Установка необходимых пакетов
RUN apk add --no-cache bash git

# Переход в рабочую директорию
WORKDIR /app

# Копирование файлов проекта из C:\programs\hack
COPY . /app

# Сборка проекта
RUN go mod tidy && go build -o main ./cmd/main.go

# Копирование статических файлов (если есть)
COPY ./static /app/static

# Запуск сервера
CMD ["/app/main"]