# syntax=docker/dockerfile:1.2

# TODO: Move that to `thegarii` repo so we can pull it?
# thegarii builder
FROM rust:bullseye as thegarii-builder
ENV CARGO_NET_GIT_FETCH_WITH_CLI=true
RUN --mount=type=cache,target=/var/cache/apk \
    --mount=type=cache,target=/home/rust/.cargo \
    rustup component add rustfmt \
    && git clone https://github.com/streamingfast/thegarii \
    && cd thegarii \
    && cargo build --release \
    && cp target/release/thegarii /home/rust/

FROM ubuntu:20.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
    apt-get -y install -y \
    ca-certificates libssl1.1 vim htop iotop sysstat \
    dstat strace lsof curl jq tzdata && \
    rm -rf /var/cache/apt /var/lib/apt/lists/*

RUN rm /etc/localtime && ln -snf /usr/share/zoneinfo/America/Montreal /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

ADD /firearweave /app/firearweave
COPY --from=thegarii-builder /home/rust/thegarii /app/thegarii

# TODO: Add back later
# COPY tools/sfeth/motd_generic /etc/
# COPY tools/sfeth/99-sfeth-generic.sh /etc/profile.d/

ENTRYPOINT ["/app/firearweave"]
