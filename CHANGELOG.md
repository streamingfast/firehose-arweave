# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). See [MAINTAINERS.md](./MAINTAINERS.md)
for instructions to keep up to date.

## v1.1.0

* Bumped all firehose dependencies (from old versions)
* Logs are now more verbose by default output to STDERR
* Removed reader-node-merge-threshold-block-age (no more merging directly in the reader node, only the merger does it now)
* New 'ready' stat in prometheus for services

## v1.0.0

#### Flags and environment variables

* Renamed the `mindreader` application to `reader`
* Renamed `common-one-blocks-store-url` to `common-one-block-store-url`
* Renamed all the `mindreader-node-*` flags to `reader-node-*`





