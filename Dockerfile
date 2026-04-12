FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
ENV GOFLAGS=-mod=mod

RUN go mod download
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM alpine:3.23 AS staging

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml

ARG ENV_FILE
COPY --from=builder /app/${ENV_FILE:-'.env'} .

EXPOSE 8080

CMD ["/app/main"]

FROM alpine:3.23 AS remote

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml

EXPOSE 8080

CMD ["/app/main"]
