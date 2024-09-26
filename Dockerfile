# Use official Golang image as build stage
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o exoplanet-service .

# Use minimal Alpine image for final stage
FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/exoplanet-service .

EXPOSE 8080

CMD ["./exoplanet-service"]

