FROM golang:1.26 AS builder

WORKDIR /src

# Зависимости отдельным слоем — кешируется между сборками.
COPY go.mod go.sum ./
RUN go mod download

# Исходники (включая пакет tests — его импортирует main.go).
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/scheduler .

# ---------- stage 2: рантайм ----------
FROM ubuntu:latest

# tzdata — чтобы time.Now() и NextDate() имели корректную таймзону.
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates tzdata \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Бинарник и фронтенд. WebDir в main.go = "web", поэтому кладём
# ровно в /app/web, чтобы http.FileServer нашёл файлы.
COPY --from=builder /out/scheduler /app/scheduler
COPY --from=builder /src/web     /app/web

# Дефолты. Все можно переопределить через `docker run -e ...`.
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

# Порт веб-сервера.
EXPOSE 7540

# Директорию для БД удобно монтировать с хоста.
VOLUME ["/data"]

ENTRYPOINT ["/app/scheduler"]