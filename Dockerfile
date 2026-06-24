# Stage 1: Builder
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/publisher ./cmd/publisher

# Stage 2: Final image (minimal)
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bin/server ./server
COPY --from=builder /app/bin/worker ./worker
COPY --from=builder /app/bin/publisher ./publisher
COPY migrations/ ./migrations/
EXPOSE 8080
CMD ["./server"]
