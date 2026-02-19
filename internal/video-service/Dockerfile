FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o video-service ../cmd/app.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .
COPY config/app.yml config/app.yml
# Optional if env is used
COPY .env .env

# Expose port (adjust if needed after checking app.go)
EXPOSE 8081

CMD ["./user-service"]
