FROM golang:1.22.5 AS build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY . .
RUN make build

FROM alpine:3.20.2
ENV DATA_SOURCE_NAME=dm://SYSDBA:SYSDBA@localhost:5236?autoCommit=true
ENV DEFAULT_METRICS=/opt/dm_exporter/cnf/default-metrics.toml
WORKDIR /opt/dm_exporter
COPY --from=build /workspace/build/dm_exporter /opt/dm_exporter/bin/dm_exporter
ADD default-metrics.toml /opt/dm_exporter/cnf/default-metrics.toml
RUN chmod 755 /opt/dm_exporter/bin/dm_exporter
EXPOSE 9161
ENTRYPOINT ["/opt/dm_exporter/bin/dm_exporter"]
