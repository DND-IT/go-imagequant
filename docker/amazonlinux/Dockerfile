FROM amazonlinux:2023 as rust
LABEL stage=rust
WORKDIR /build
RUN dnf update -y  && \
    dnf install gcc openssl-devel git -y
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh  -s -- -y
RUN /bin/sh -c 'export PATH=$PATH:~/.cargo/bin/ && \
    cargo install cargo-c'
RUN git clone https://github.com/ImageOptim/libimagequant.git --branch 4.2.0
RUN /bin/sh -c 'export PATH=$PATH:~/.cargo/bin/ && \
    cd libimagequant/imagequant-sys && \
    CFLAGS=-mno-outline-atomics cargo cinstall --destdir=. '

# compile
FROM amazonlinux:2023 as builder
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
LABEL stage=builder
RUN dnf update -y && \
    dnf install -y \
    gcc g++ \
    golang \
    make \
    openssl \
    git \
    libwebp-devel
WORKDIR /build
COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/lib/ /usr/lib64
COPY . .
RUN /bin/sh -c 'make test'
RUN /bin/sh -c 'make build'
# build image
FROM amazonlinux:2023
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
ENV LANG en_US.UTF-8
ENV TERM linux
# add needed packages
RUN dnf update -y
RUN dnf install -y libwebp &&  \
    mkdir -p /opt/bin
COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/lib/ /usr/lib64
WORKDIR /opt/bin
COPY --from=builder /build/cmd/go-imagequant .
ENTRYPOINT ["./go-imagequant"]