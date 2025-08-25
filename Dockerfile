# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./

# Install gcc and musl-dev for CGO
RUN apk add --no-cache gcc musl-dev

RUN go mod download
COPY . .

# Enable CGO for sqlite3
ENV CGO_ENABLED=1
RUN go build -o purr .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/purr .
COPY pixel.gif .
# Optional: copy an empty SQLite DB or let purr create it at runtime
EXPOSE 8080
CMD ["./purr"]
