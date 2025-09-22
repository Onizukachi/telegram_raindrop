# Build Stage
FROM golang:alpine3.22 as builder

WORKDIR /builder

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0  go build -o /app ./cmd/bot

# Final Stage
FROM alpine:3.20

COPY --from=builder /app /app
COPY messages.yml ./
COPY ./.env ./
COPY ./internal/storage/migrations ./internal/storage/migrations

CMD ["/app"]