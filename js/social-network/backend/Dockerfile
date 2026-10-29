FROM golang:1.22-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o social-network .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/social-network .
COPY --from=builder /app/pkg/db/migrations/sqlite ./pkg/db/migrations/sqlite

EXPOSE 8080

CMD ["./social-network"]
