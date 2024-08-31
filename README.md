# otel-exporter-demo-go

This code demonstrates how to use OpenTelemetry to send metrics from a Golang application to an OpenTelemetry Collector using HTTP protocol.

## Features

- Initializes OpenTelemetry SDK
- Sets up a metrics provider
- Creates and increments a counter metric
- Creates and observes a gauge metric with random values
- Sends metrics periodically (every 1 second)

## Configuration

The application is configured to send metrics to `localhost:4318` by default. If your OpenTelemetry Collector is running on a different address, update the following line in the `initProvider` function:

```go
otlpmetrichttp.WithEndpoint("localhost:4318"),
```

## Usage

To run the application:

```
go run main.go
```

The application will start sending metrics to the configured OpenTelemetry Collector. You should see output in your collector logs indicating that it's receiving metrics.

## Example using OpenTelemetry Collector and Prometheus

Example of a configuration that sends metrics from otel-exporter-demo-go to OpenTelemetry Collector and stores them in Prometheus.

### Set up OpenTelemetry Collector

https://github.com/open-telemetry/opentelemetry-collector-releases

```
$ git clone git@github.com:open-telemetry/opentelemetry-collector-releases.git
$ cd opentelemetry-collector-releases
$ make ocb
Installing ocb (linux/amd64) at /home/hoge/bin
$ ocb --config distributions/otelcol/manifest.yaml
2024-08-31T06:33:50.720+0900    INFO    internal/command.go:125 OpenTelemetry Collector Builder{"version": "v0.108.1"}
2024-08-31T06:33:50.722+0900    INFO    internal/command.go:161 Using config file       {"path": "distributions/otelcol/manifest.yaml"}
2024-08-31T06:33:50.722+0900    INFO    builder/config.go:142   Using go        {"go-executable": "/usr/bin/go"}
2024-08-31T06:33:50.723+0900    INFO    builder/main.go:101     Sources created {"path": "./_build"}
2024-08-31T06:33:51.106+0900    INFO    builder/main.go:192     Getting go modules
2024-08-31T06:33:51.184+0900    INFO    builder/main.go:112     Compiling
2024-08-31T06:33:55.006+0900    INFO    builder/main.go:131     Compiled        {"binary": "./_build/otelcol"}
$ _build/otelcol --version
otelcol version 0.108.1
```

```
$ cat << EOF config.yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:

exporters:
  debug:
    verbosity: detailed

  prometheusremotewrite:
    endpoint: "http://localhost:9090/api/v1/write"
    tls:
      insecure: true

service:

  pipelines:

    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug, prometheusremotewrite]
EOF
```

```
$ _build/otelcol --config config.yaml
2024-08-31T10:40:46.655+0900    info    service@v0.108.1/service.go:178 Setting up own telemetry...
2024-08-31T10:40:46.655+0900    info    service@v0.108.1/telemetry.go:98        Serving metrics{"address": ":8888", "metrics level": "Normal"}
2024-08-31T10:40:46.656+0900    info    builders/builders.go:26 Development component. May change in the future.        {"kind": "exporter", "data_type": "metrics", "name": "debug"}
2024-08-31T10:40:46.657+0900    info    service@v0.108.1/service.go:263 Starting otelcol...    {"Version": "0.108.1", "NumCPU": 8}
2024-08-31T10:40:46.657+0900    info    extensions/extensions.go:38     Starting extensions...
2024-08-31T10:40:46.657+0900    info    otlpreceiver@v0.108.1/otlp.go:153       Starting HTTP server    {"kind": "receiver", "name": "otlp", "data_type": "metrics", "endpoint": "0.0.0.0:4318"}
2024-08-31T10:40:46.657+0900    info    service@v0.108.1/service.go:289 Everything is ready. Begin running and processing data.
2024-08-31T10:40:46.657+0900    info    localhostgate/featuregate.go:63 The default endpoints for all servers in components have changed to use localhost instead of 0.0.0.0. Disable the feature gate to temporarily revert to the previous default.   {"feature gate ID": "component.UseLocalHostAsDefaultHost"}
```

### Set up Prometheus

https://github.com/prometheus/prometheus/releases/tag/v2.54.1

```
$ ./prometheus --version
prometheus, version 2.54.1 (branch: HEAD, revision: e6cfa720fbe6280153fab13090a483dbd40bece3)
  build user:       root@812ffd741951
  build date:       20240827-10:56:41
  go version:       go1.22.6
  platform:         linux/amd64
  tags:             netgo,builtinassets,stringlabels
```

```
$ ./prometheus --web.enable-remote-write-receiver
ts=2024-08-31T01:52:11.760Z caller=main.go:601 level=info msg="No time or size retention was set so using the default time retention" duration=15d
ts=2024-08-31T01:52:11.760Z caller=main.go:645 level=info msg="Starting Prometheus Server" mode=server version="(version=2.54.1, branch=HEAD, revision=e6cfa720fbe6280153fab13090a483dbd40bece3)"
ts=2024-08-31T01:52:11.760Z caller=main.go:650 level=info build_context="(go=go1.22.6, platform=linux/amd64, user=root@812ffd741951, date=20240827-10:56:41, tags=netgo,builtinassets,stringlabels)"
ts=2024-08-31T01:52:11.760Z caller=main.go:651 level=info host_details="(Linux 6.9.12-amd64 #1 SMP PREEMPT_DYNAMIC Debian 6.9.12-1 (2024-07-27) x86_64 fumino (none))"
ts=2024-08-31T01:52:11.760Z caller=main.go:652 level=info fd_limits="(soft=1073741816, hard=1073741816)"
ts=2024-08-31T01:52:11.760Z caller=main.go:653 level=info vm_limits="(soft=unlimited, hard=unlimited)"
ts=2024-08-31T01:52:11.762Z caller=web.go:571 level=info component=web msg="Start listening for connections" address=0.0.0.0:9090
ts=2024-08-31T01:52:11.763Z caller=main.go:1160 level=info msg="Starting TSDB ..."
ts=2024-08-31T01:52:11.767Z caller=tls_config.go:313 level=info component=web msg="Listening on" address=[::]:9090
ts=2024-08-31T01:52:11.767Z caller=tls_config.go:316 level=info component=web msg="TLS is disabled." http2=false address=[::]:9090
ts=2024-08-31T01:52:11.768Z caller=head.go:626 level=info component=tsdb msg="Replaying on-disk memory mappable chunks if any"
ts=2024-08-31T01:52:11.768Z caller=head.go:713 level=info component=tsdb msg="On-disk memory mappable chunks replay completed" duration=4.145µs
ts=2024-08-31T01:52:11.768Z caller=head.go:721 level=info component=tsdb msg="Replaying WAL, this may take a while"
ts=2024-08-31T01:52:11.770Z caller=head.go:793 level=info component=tsdb msg="WAL segment loaded" segment=0 maxSegment=5
ts=2024-08-31T01:52:11.771Z caller=head.go:793 level=info component=tsdb msg="WAL segment loaded" segment=1 maxSegment=5
ts=2024-08-31T01:52:11.773Z caller=head.go:793 level=info component=tsdb msg="WAL segment loaded" segment=2 maxSegment=5
ts=2024-08-31T01:52:11.774Z caller=head.go:793 level=info component=tsdb msg="WAL segment loaded" segment=3 maxSegment=5
ts=2024-08-31T01:52:11.775Z caller=head.go:793 level=info component=tsdb msg="WAL segment loaded" segment=4 maxSegment=5
ts=2024-08-31T01:52:11.775Z caller=head.go:793 level=info component=tsdb msg="WAL segment loaded" segment=5 maxSegment=5
ts=2024-08-31T01:52:11.775Z caller=head.go:830 level=info component=tsdb msg="WAL replay completed" checkpoint_replay_duration=44.429µs wal_replay_duration=6.992432ms wbl_replay_duration=140ns chunk_snapshot_load_duration=0s mmap_chunk_replay_duration=4.145µs total_replay_duration=7.064704ms
ts=2024-08-31T01:52:11.778Z caller=main.go:1181 level=info fs_type=EXT4_SUPER_MAGIC
ts=2024-08-31T01:52:11.778Z caller=main.go:1184 level=info msg="TSDB started"
ts=2024-08-31T01:52:11.778Z caller=main.go:1367 level=info msg="Loading configuration file" filename=prometheus.yml
ts=2024-08-31T01:52:11.779Z caller=main.go:1404 level=info msg="updated GOGC" old=100 new=75
ts=2024-08-31T01:52:11.779Z caller=main.go:1415 level=info msg="Completed loading of configuration file" filename=prometheus.yml totalDuration=458.915µs db_storage=777ns remote_storage=1.134µs web_handler=615ns query_engine=6.464µs scrape=159.723µs scrape_sd=24.351µs notify=17.913µs notify_sd=6.38µs rules=1.031µs tracing=4.634µs
ts=2024-08-31T01:52:11.779Z caller=main.go:1145 level=info msg="Server is ready to receive web requests."
ts=2024-08-31T01:52:11.779Z caller=manager.go:164 level=info component="rule manager" msg="Starting rule manager..."
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
