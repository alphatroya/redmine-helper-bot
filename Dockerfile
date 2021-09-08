FROM golang:1.17 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.9

RUN adduser -D -g '' botuser
USER botuser

WORKDIR /home/botuser
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
