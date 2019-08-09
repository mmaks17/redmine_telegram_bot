# FROM golang:1.12
# WORKDIR /usr/bin

# COPY ./BotIntraservice /

# CMD ["/BotIntraservice"]

FROM golang:latest AS build

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0
RUN go get && go build -o app main.go

FROM alpine:latest
WORKDIR /app
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=build /go/src/app/app /app/

CMD ["/app/app"]


