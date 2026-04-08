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
ARG USER_DB_HOST
ARG USER_DB_PORT
ARG USER_API_PORT
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD

# Set environment variables in the container
ENV USER_DB_NAME=${USER_DB_NAME}
ENV USER_DB_HOST=${USER_DB_HOST}
ENV USER_DB_PORT=${USER_DB_PORT}
ENV USER_API_PORT=${USER_API_PORT}
ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

COPY --from=builder /app/user-service .

CMD ["./user-service"]
