FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o user-service cmd/app.go

FROM alpine:latest

WORKDIR /app

# Load environment variables from .env
ARG USER_DB_NAME
ARG USER_DB_USER
ARG USER_DB_PASSWORD
ARG USER_DB_HOST
ARG USER_DB_PORT
ARG USER_API_PORT

# Set environment variables in the container
ENV USER_DB_NAME=${USER_DB_NAME}
ENV USER_DB_USER=${USER_DB_USER}
ENV USER_DB_PASSWORD=${USER_DB_PASSWORD}
ENV USER_DB_HOST=${USER_DB_HOST}
ENV USER_DB_PORT=${USER_DB_PORT}
ENV USER_API_PORT=${USER_API_PORT}

COPY --from=builder /app/user-service .

CMD ["./user-service"]
