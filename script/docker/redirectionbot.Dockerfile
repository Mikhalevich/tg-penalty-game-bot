FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/redirectionbot cmd/redirectionbot/main.go

FROM alpine:3.23

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/redirectionbot /app/redirectionbot
#COPY --from=builder /app/config/config-redirectionbot.yaml /app/config/config-redirectionbot.yaml

ENTRYPOINT ["./redirectionbot", "-config", "./config/config-redirectionbot.yaml"]
