# =========================
# BUILD STAGE
# =========================
FROM golang:1.25-alpine AS builder

# Install CGO deps for SQLite
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy Go module files first (cache-friendly)
COPY go.mod go.sum ./
RUN go mod download

# Copy entire source
COPY . .

# Build the forum application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o forum ./cmd/app

# =========================
# RUNTIME STAGE
# =========================
FROM alpine:3.19

# Install SQLite runtime libraries
RUN apk add --no-cache sqlite-libs

# Metadata labels
LABEL org.opencontainers.image.title="Go Forum"
LABEL org.opencontainers.image.description="Forum web application built with Go, HTML, and SQLite"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.authors="team@forum.dev"

WORKDIR /app

# Copy binary
COPY --from=builder /app/forum .

# Copy frontend & templates
COPY web ./web

# Copy migrations & seed data
COPY migrations ./migrations
COPY internal/db ./internal/db

# Create directory for SQLite database
RUN mkdir -p /data

# Expose web port
EXPOSE 8080

# SQLite DB path (used by app)
ENV DB_PATH=/data/app.db

# Start app
CMD ["./forum"]