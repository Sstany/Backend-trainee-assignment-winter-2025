FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o core /app/cmd/core/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/core /app/core

EXPOSE 8080

CMD [ "./core" ]
