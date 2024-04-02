FROM buildpack-deps:jammy as rust

WORKDIR /build

RUN apt-get update -y  && \
    apt-get install -y --no-install-recommends git libssl-dev
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh  -s -- -y
RUN /bin/bash -c 'source "$HOME/.cargo/env" && \
    cargo install cargo-c'
RUN git clone https://github.com/ImageOptim/libimagequant.git --branch 4.2.2
RUN /bin/bash -c 'source "$HOME/.cargo/env" && \
    cd libimagequant/imagequant-sys && \
    CFLAGS=-mno-outline-atomics cargo cinstall --destdir=. '

FROM scratch as artifact
COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/ /lib/ubuntu/22.04/