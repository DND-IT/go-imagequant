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

FROM scratch as artifact
COPY --from=rust /build/libimagequant/imagequant-sys/usr/local /lib/alpine/3.17/