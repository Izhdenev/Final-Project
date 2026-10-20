# ---------- stage 1: сборка ----------
FROM golang:1.26 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/scheduler .

# ---------- stage 2: рантайм ----------
FROM alpine:latest

# ca-certificates — если понадобится HTTPS наружу,
# tzdata — чтобы time.Now() и NextDate() знали локальную таймзону.
RUN apk add --no-cache ca-certificates tzdata

# Непривилегированный пользователь — хорошая практика.
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/scheduler /app/scheduler
COPY --from=builder /src/web     /app/web

RUN mkdir -p /data && chown -R app:app /app /data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

VOLUME ["/data"]

USER app

ENTRYPOINT ["/app/scheduler"]