# Gunakan image Go sebagai base image
# FROM golang:1.20-alpine AS builder
FROM golang:1.21 AS builder

# Set lingkungan kerja di dalam container
WORKDIR /app

# Copy go.mod dan go.sum untuk menginstall dependencies
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy seluruh source code ke dalam container
COPY . .

# Build aplikasi Go
RUN go build -o main .

# Stage kedua untuk mengurangi ukuran image (menggunakan Alpine sebagai base image)
FROM alpine:latest

# Set lingkungan kerja
WORKDIR /root/

# Copy file binary dari builder stage
COPY --from=builder /app/main .

# Eksekusi binary ketika container dijalankan
CMD ["./main"]
