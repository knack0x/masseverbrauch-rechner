FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY web/ web/
COPY VERSION .
RUN go build -o web-server ./web

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/web-server .
COPY --from=builder /app/VERSION .
EXPOSE 80
ENV PORT=80
CMD ["./web-server"]
