FROM golang:1.13-alpine AS builder

RUN apk update && apk add git gcc linux-headers musl-dev
WORKDIR /go
ADD faucet.go .
ADD faucet.html .
RUN go get github.com/ethereum/go-ethereum/rpc
# `go get ./...` is not working for me
RUN go build faucet.go
  
FROM alpine:latest AS production

COPY --from=builder /go/faucet /usr/local/bin/
COPY --from=builder /go/faucet.html /var/www/

EXPOSE 8088
ENTRYPOINT ["faucet"]