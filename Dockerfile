FROM golang:1.13-alpine AS dependencies
WORKDIR /go
RUN apk update && apk add git gcc linux-headers musl-dev
# `go get ./...` is not working for me
RUN go get github.com/ethereum/go-ethereum/rpc

FROM offcode/lightfaucet:dependencies AS builder
WORKDIR /go
ADD faucet.go .
ADD faucet.html .
RUN apk add git
RUN go get github.com/didip/tollbooth
RUN go get github.com/didip/tollbooth/limiter
RUN go build faucet.go
  
FROM alpine:latest AS production
COPY --from=builder /go/faucet /usr/local/bin/
COPY --from=builder /go/faucet.html /var/www/
EXPOSE 8088
ENTRYPOINT ["faucet"]