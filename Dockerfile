FROM golang:1.19-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./api-gateway ./cmd

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/api-gateway .

COPY ./translation ./translation
COPY ./config.yaml .

EXPOSE 8080

ENTRYPOINT ["/api-gateway"]