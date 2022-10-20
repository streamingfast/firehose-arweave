FROM ubuntu:20.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
    apt-get -y install -y \
    ca-certificates libssl1.1 vim htop iotop sysstat \
    dstat strace lsof curl jq tzdata && \
    rm -rf /var/cache/apt /var/lib/apt/lists/*

RUN rm /etc/localtime && ln -snf /usr/share/zoneinfo/America/Montreal /etc/localtime && dpkg-reconfigure -f noninteractive tzdata
RUN mkdir -p /app/ && curl -Lo /app/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.12/grpc_health_probe-linux-amd64 && chmod +x /app/grpc_health_probe
RUN curl -Lo /app/thegarii https://github.com/streamingfast/thegarii/releases/download/v0.2.1/thegarii-x86_64-unknown-linux-gnu && chmod +x /app/thegarii

ADD /firearweave /app/firearweave

COPY tools/firearweave/motd_generic /etc/motd
COPY tools/firearweave/99-firehose.sh /etc/profile.d/

# On SSH connection, /root/.bashrc is invoked which invokes '/root/.bash_aliases' if existing,
# so we hijack the file to "execute" our specialized bash script
RUN echo ". /etc/profile.d/99-firehose.sh" > /root/.bash_aliases

ENV PATH "$PATH:/app"

ENTRYPOINT ["/app/firearweave"]
