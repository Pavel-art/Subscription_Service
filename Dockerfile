########################
# Build stage
########################
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o subscription-service ./cmd

########################
# Runtime stage
########################
FROM alpine:3.19

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /app/subscription-service .
COPY --from=builder /app/migrations ./migrations

ENV GIN_MODE=release
ENV HTTP_PORT=8081

EXPOSE 8081

USER appuser

ENTRYPOINT ["./subscription-service"]