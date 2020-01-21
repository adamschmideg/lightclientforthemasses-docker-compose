FROM golang:1.13-alpine AS builder
RUN apk update && apk add git gcc linux-headers musl-dev
ADD . /lightfaucet
RUN cd /lightfaucet && go build

FROM alpine:latest AS production
COPY --from=builder /lightfaucet/faucet /usr/local/bin/
COPY --from=builder /lightfaucet/faucet.html /var/www/
EXPOSE 8088
ENTRYPOINT ["faucet"]