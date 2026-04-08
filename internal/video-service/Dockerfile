FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o video-service cmd/app.go

FROM alpine:latest

WORKDIR /app

# Load environment variables from .env
ARG VIDEO_DB_NAME
ARG VIDEO_DB_HOST
ARG VIDEO_DB_PORT
ARG VIDEO_API_PORT
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD

# Set environment variables in the container
ENV VIDEO_DB_NAME=${VIDEO_DB_NAME}
ENV VIDEO_DB_HOST=${VIDEO_DB_HOST}
ENV VIDEO_DB_PORT=${VIDEO_DB_PORT}
ENV VIDEO_API_PORT=${VIDEO_API_PORT}
ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

COPY --from=builder /app/video-service .

CMD ["./video-service"]
