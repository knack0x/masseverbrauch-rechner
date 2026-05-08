FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
COPY main.go .
COPY web/ web/
COPY VERSION .
RUN go build -o web-server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/web-server .
COPY --from=builder /app/VERSION .
EXPOSE 80
ENV PORT=80
CMD ["./web-server"]
