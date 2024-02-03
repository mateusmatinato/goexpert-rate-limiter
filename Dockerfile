FROM golang:1.21 AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rate-limiter cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/rate-limiter .
COPY ./configs/config.env ./configs/config.env
ENTRYPOINT ["./rate-limiter"]