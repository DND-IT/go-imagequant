FROM rust:alpine3.15 as rust
LABEL stage=rust
WORKDIR /build
COPY ./docker/alpine-libimagequant/enable-world-repos.sh /root/
RUN /root/enable-world-repos.sh
# RUN PATH="/usr/local/cargo/bin/:$PATH"
# RUN export PATH
RUN /bin/sh -c 'apk update && \
                apk upgrade && \
                apk add bash curl bind-tools sudo build-base git musl-dev libressl-dev cargo-c && \
                git clone https://github.com/ImageOptim/libimagequant.git --branch 4.0.0'

# RUN /bin/sh -c 'cd libimagequant &&  \
#    CFLAGS=-mno-outline-atomics cargo install cargo-c &&  \
#    ls -l /usr/local/cargo/bin'

#RUN /bin/sh -c  'cd libimagequant/imagequant-sys && \
#    ls -l . && cargo cinstall --destdir=/build/usr/local'
#    # RUN /bin/sh -c 'cd libimagequant/imagequant-sys && \
##                 CFLAGS=-mno-outline-atomics cargo cinstall --destdir=. '
#
#
## build image, based on alpine:3.15
#FROM alpine:3.15
#LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
#COPY --from=rust /build/usr/local/ /usr/local/
