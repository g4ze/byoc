# Build stage
FROM golang:1.22.5 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM golang:1.22.5

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 2001

CMD ["./main"]