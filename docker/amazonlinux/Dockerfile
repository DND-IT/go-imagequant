FROM amazonlinux:2 as rust
LABEL stage=rust
WORKDIR /build
RUN yum update -y  && \
    yum install gcc openssl-devel git -y
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh  -s -- -y
RUN export PATH=$PATH:~/.cargo/bin/ && \
    cargo install cargo-c
RUN git clone https://github.com/ImageOptim/libimagequant.git --branch 4.2.0
RUN export PATH=$PATH:~/.cargo/bin/ && \
    cd libimagequant/imagequant-sys && \
    CFLAGS=-mno-outline-atomics cargo cinstall --destdir=.

# compile libwebp
FROM amazonlinux:2 AS webpbuilder
LABEL maintainer="Christian Juerges (christian.juerges@tx.group) and others"
LABEL stage=webpbuilder
WORKDIR /build
RUN yum update -y && yum -y group install "Development Tools" && \
    yum install -y libjpeg-turbo-devel libpng-devel libtiff-devel giflib-devel
RUN git clone --depth 1 --branch v1.3.1 https://chromium.googlesource.com/webm/libwebp
RUN cd libwebp && \
    ./autogen.sh && ./configure && make && make install
RUN /bin/ls -l /

# compile
FROM amazonlinux:2 as builder
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
LABEL stage=builder
RUN yum update -y && \
    yum install -y \
    gcc g++ \
    make \
    wget \
    tar \
    openssl \
    git
WORKDIR /build

RUN wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz

COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/lib/ /usr/lib64
COPY --from=webpbuilder /usr/local/include/webp/ /usr/local/include/webp/
COPY --from=webpbuilder /usr/local/lib/ /usr/lib64/
COPY --from=webpbuilder /usr/local/lib/pkgconfig/ /usr/local/lib/pkgconfig/

COPY . .
RUN PATH=$PATH:/usr/local/go/bin make test
RUN PATH=$PATH:/usr/local/go/bin make build

# build image
FROM amazonlinux:2
LABEL maintainer="Christian Juerges, christian.juerges@tx.group"
ENV LANG en_US.UTF-8
ENV TERM linux
# add needed packages
RUN yum update -y && \
    yum clean all && \
    rm -rf /var/cache/yum
RUN mkdir -p /opt/bin
COPY --from=rust /build/libimagequant/imagequant-sys/usr/local/lib/ /usr/lib64
COPY --from=webpbuilder /usr/local/lib/ /usr/lib64/
WORKDIR /opt/bin
COPY --from=builder /build/cmd/go-imagequant .
ENTRYPOINT ["./go-imagequant"]