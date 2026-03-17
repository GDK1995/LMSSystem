FROM golang:1.25-alpine AS build

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ---- runtime image ----
FROM alpine:latest

WORKDIR /app

# копируем бинарник
COPY --from=build /app/main .

# если есть миграции
# COPY --from=build /app/db/migrations ./db/migrations

EXPOSE 8082

CMD ["./main"]