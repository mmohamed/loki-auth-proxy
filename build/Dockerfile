FROM golang:1.12-alpine as builder

ARG VERSION=dev

ENV GO111MODULE=on
ENV CGO_ENABLED=0

ENV SRC_PATH=github.com/medinvention/loki-auth-proxy

RUN apk add  -U --no-cache git ca-certificates && \
    mkdir -p /go/src/${SRC_PATH}

WORKDIR /go/src/${SRC_PATH}

COPY go.mod go.mod
COPY go.sum go.sum
COPY src src

RUN cd src/cmd && \
    go build -ldflags="-X 'main.version=${VERSION}'" \ 
    -o /go/src/${SRC_PATH}/build/proxy

FROM scratch

ENV SRC_PATH=github.com/medinvention/loki-auth-proxy

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/${SRC_PATH}/build/proxy /proxy

ENTRYPOINT [ "/proxy" ]