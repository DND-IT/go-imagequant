FROM alpine:3.17  as rust
LABEL stage=rust
WORKDIR /build
COPY ./docker/alpine/enable-world-repos.sh /root/
RUN /root/enable-world-repos.sh
RUN apk update && \
    apk upgrade && \
    apk add rust cargo cargo-c bash curl bind-tools sudo build-base git musl-dev libressl-dev && \
    git clone https://github.com/ImageOptim/libimagequant.git --branch 4.0.1

RUN /bin/sh -c 'cd libimagequant/imagequant-sys && \
    cargo cinstall --destdir=. '

# compile
FROM golang:1.18.8-alpine3.17 as builder
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
LABEL stage=builder

RUN set -eux; \
    apk update && apk upgrade && apk add bash \
    gcc g++ \
    make \
    musl-dev \
    openssl \
    git \
    libwebp-dev

COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/ /usr/local/

WORKDIR /build
COPY . .
RUN /bin/sh -c 'make test'
RUN /bin/sh -c 'make build'

# build image, based on alpine:3.15
FROM alpine:3.17
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"

ENV LANG en_US.UTF-8
ENV TERM linux

# add needed packages
RUN /bin/sh -c 'apk update && \
    apk upgrade && \
    apk add bash curl bind-tools sudo libwebp libgcc bcc-tools bcc-doc && \
    mkdir -p /opt/bin'


COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/ /usr/local/

WORKDIR /opt/bin
COPY --from=builder /build/cmd/go-imagequant .
ENTRYPOINT ["./go-imagequant"]

