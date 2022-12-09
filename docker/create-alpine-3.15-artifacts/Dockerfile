FROM rust:alpine3.15 as rust
LABEL stage=rust
WORKDIR /build
RUN apk update && \
    apk upgrade && \
    apk add bash curl bind-tools sudo build-base git musl-dev libressl-dev
RUN git clone https://github.com/ImageOptim/libimagequant.git --branch 4.0.0
RUN /bin/sh -c 'export PATH=$PATH:~/.cargo/bin/ && \
    cd libimagequant/imagequant-sys && \
    CFLAGS=-mno-outline-atomics cargo cinstall --destdir=. '

FROM scratch as artifact
COPY --from=rust /build/libimagequant/imagequant-sys/usr/local /lib/alpine/3.15/