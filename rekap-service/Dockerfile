# STAGE 1: Builder
FROM golang:1.26-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o rekap-app .

# STAGE 2: Runner
FROM alpine:3.21
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN adduser -D appuser
COPY --from=builder /app/rekap-app .
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./rekap-app"]