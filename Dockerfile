ARG GO_VERSION=1.14.3
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /
ENV CGO_ENABLED=0
COPY ./go.mod ./go.sum ./main.go ./ 
RUN go build

CMD ./soil-moisture-ws