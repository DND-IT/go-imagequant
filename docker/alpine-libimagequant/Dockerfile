FROM rust:alpine3.15 as rust
LABEL stage=rust
WORKDIR /build
COPY ./docker/alpine-libimagequant/enable-world-repos.sh /root/
RUN /root/enable-world-repos.sh
RUN apk update && \
                apk upgrade && \
                apk add bash curl bind-tools sudo build-base git musl-dev libressl-dev cargo-c && \
                git clone https://github.com/ImageOptim/libimagequant.git --branch 4.0.0

RUN /bin/sh -c 'cd libimagequant/imagequant-sys && \
                CFLAGS=-mno-outline-atomics cargo cinstall --destdir=. '


# add to golang base image
FROM golang:1.16.12-alpine3.15 as builder
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/ /usr/local/
