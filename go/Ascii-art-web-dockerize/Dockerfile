FROM golang:1.23-alpine AS hassan

WORKDIR /omran

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o server

EXPOSE 8080

CMD ["./server"]
