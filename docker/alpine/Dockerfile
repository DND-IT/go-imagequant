# compile
FROM golang:1.21.8-alpine3.18 as builder
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
LABEL stage=builder

RUN set -eux; \
    apk update && apk upgrade && apk add bash \
    gcc g++ \
    make \
    musl-dev \
    openssl \
    git \
    libwebp-dev \
    libimagequant-dev

WORKDIR /build
COPY . .
RUN /bin/sh -c 'make test'
RUN /bin/sh -c 'make build'

# build image, based on alpine
FROM alpine:3.18
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"

ENV LANG en_US.UTF-8
ENV TERM linux

# add needed packages
RUN /bin/sh -c 'apk update && \
    apk upgrade && \
    apk add bash curl bind-tools sudo libwebp libimagequant libgcc bcc-tools bcc-doc  && \
    mkdir -p /opt/bin'



WORKDIR /opt/bin
COPY --from=builder /build/cmd/go-imagequant .
ENTRYPOINT ["./go-imagequant"]

