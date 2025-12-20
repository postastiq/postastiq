FROM golang:1.21-alpine AS builder

# Install build dependencies including gcc for CGO (required by sqlite3)
RUN apk add --no-cache git ca-certificates gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o blog .

# Final stage
FROM alpine:latest

# Install ca-certificates and sqlite libraries
RUN apk --no-cache add ca-certificates tzdata sqlite-libs

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/blog .

# Create a volume for the database
VOLUME /app/data

# Expose port
EXPOSE 8080

# Set default database path
ENV DB_PATH=/app/data/blog.db

# Run the application
CMD ["./blog"]
