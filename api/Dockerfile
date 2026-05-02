FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY api/ .
COPY VERSION .
RUN go build -o api main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/VERSION .
EXPOSE 8080
CMD ["./api"]
