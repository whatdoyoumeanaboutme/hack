FROM golang:1.20-alpine

# Установка необходимых пакетов
RUN apk add --no-cache bash git

# Переход в рабочую директорию
WORKDIR /app

# Копирование файлов проекта
COPY . /app

# Сборка проекта
RUN go mod tidy && go build -o main ./cmd/main.go

# Копирование файлов вебсайта (если есть)
COPY ./css /app/css
COPY ./js /app/js

# Запуск сервера
CMD ["/app/main"]