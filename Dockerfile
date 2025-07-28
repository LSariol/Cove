# import base image
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cove ./cmd/cove


FROM alpine:latest
COPY --from=builder /app/cove /cove
EXPOSE 8081
CMD ["/cove"]