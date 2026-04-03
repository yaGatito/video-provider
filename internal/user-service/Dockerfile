FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o user-service cmd/user-service/app.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .
COPY config/user-service.yml /app/config/user-service.yml

# TODO: consider using env vars instead of .env file for better Docker practices
COPY .env .env

CMD ["./user-service"]

#  TODO: run go natively [x: works expected] and through docker and see difference on why it cant connect to the DB
