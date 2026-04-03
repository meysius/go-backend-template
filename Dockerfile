# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app main.go

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl \
    && curl -sSf https://atlasgo.sh | sh

COPY --from=builder /app /app
COPY migrations/ /migrations/

EXPOSE 3000

ENTRYPOINT ["/app"]
