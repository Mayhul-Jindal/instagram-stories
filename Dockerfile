# builder stage
FROM golang:1.21.5-alpine3.19 as builder

WORKDIR /app

# these layers can be reused because of caching mechanism
COPY go.mod go.sum ./
RUN go mod download
RUN apk add curl
RUN  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

COPY . .
RUN go mod tidy
RUN go build -o instgram-stories-service .
RUN chmod +x instgram-stories-service

# Final stage
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/instgram-stories-service .
COPY --from=builder /app/migrate.linux-amd64 /bin/migrate
COPY postgres/migration ./postgres/migration

EXPOSE 3000

CMD migrate -path ./postgres/migration -database "postgresql://postgres:postgres@postgres:5432/instagram-stories?sslmode=disable" -verbose up && ./instgram-stories-service