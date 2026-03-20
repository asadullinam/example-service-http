FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/example-service ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /bin/example-service /example-service

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/example-service"]
