FROM golang:1.21-alpine3.19 AS builder
WORKDIR /app
COPY . /app

RUN go build -o main main.go

#build small image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["/app/main"]