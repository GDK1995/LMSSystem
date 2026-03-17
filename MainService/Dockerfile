FROM golang:1.25-alpine AS build

# рабочая директория
WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# -------- runtime stage --------
FROM alpine:latest

# рабочая директория
WORKDIR /app

# копируем бинарник
COPY --from=build /app/main .

# если есть миграции
# COPY --from=build /app/db/migrations ./db/migrations

# порт сервиса
EXPOSE 8083

# запуск
CMD ["./main"]