FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /app
COPY . .

ENV GO111MODULE=on
ENV GOFLAGS=-mod=mod

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 9000

CMD ["/app/main"]
