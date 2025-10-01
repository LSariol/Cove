# -- Build --
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cove ./cmd/cove

# -- Final --
FROM alpine:latest
WORKDIR /app

RUN mkdir -p /app/cove

COPY --from=builder /app/cove /cove

EXPOSE 2100
CMD ["/cove"]