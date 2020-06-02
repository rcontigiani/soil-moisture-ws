FROM golang:alpine as builder
WORKDIR /app
ADD . /app
RUN cd /app && go build -o goapp

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /app/goapp /app

EXPOSE 8080
ENTRYPOINT ./goapp