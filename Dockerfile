FROM golang:1.17-buster AS builder
WORKDIR /src
COPY . .
RUN apt-get update && apt-get -y install cmake libssl-dev
RUN ./scripts/install_libgit2.sh
RUN go mod download
RUN make

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/.build/server /server
ENTRYPOINT [ "/server" ]

EXPOSE 3306
