# ---------- Build Stage ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git (needed sometimes for modules)
RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download

COPY . .

# Build binary (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o flowgo

# ---------- Runtime Stage ----------
FROM alpine:3.20

WORKDIR /app

# Add non-root user (important for production)
RUN adduser -D appuser

COPY --from=builder /app/flowgo .
COPY example.json .

USER appuser

ENTRYPOINT ["./flowgo"]
CMD ["--load", "example.json"]