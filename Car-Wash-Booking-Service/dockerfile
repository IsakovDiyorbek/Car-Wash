FROM golang:1.22.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o root .

FROM alpine:latest
WORKDIR /app
COPY .env .
COPY --from=builder /app/root .

EXPOSE 1212

CMD ["./root"]
    