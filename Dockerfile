# Сначала используемgolang:1.24 как builder
FROM golang:1.24-alpine AS builder

# Установка необходимых пакетов для компиляции
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum* ./

# Если замены локальных пакетов (require ... => ./internal), копируем их
COPY internal/ ./internal/
COPY cmd/ ./cmd/

# Загружаем зависимости
RUN go mod download

# Копируем остальные исходные файлы
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=1 GOOS=linux go build -o /glavredusgo ./cmd/main.go

# Финальный образ - минимальный
FROM alpine:latest

# Установка необходимых библиотек для SQLite
RUN apk add --no-cache ca-certificates tzdata musl

WORKDIR /app

# Создаем директории для данных
RUN mkdir -p data history.bleve

# Копируем бинарник из builder
COPY --from=builder /glavredusgo /app/glavredusgo

# Копируем стартовые данные (history.bleve)
COPY --from=builder /app/history.bleve /app/history.bleve

EXPOSE 8080

# Запуск приложения
CMD ["./glavredusgo"]