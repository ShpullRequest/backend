FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY .env .

COPY ./docs /app/docs

CMD ["sh", "-c", "while ! nc -z postgres 5432; do sleep 2; done; ./app"]
