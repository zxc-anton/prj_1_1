FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o /go-app ./cmd/main.go

FROM alpine:3.20

WORKDIR /root

COPY --from=builder /go-app /go-app
COPY --from=builder /app/config.yaml /root/config.yaml

EXPOSE 8000

CMD [ "/go-app" ]