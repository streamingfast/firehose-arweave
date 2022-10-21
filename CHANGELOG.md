# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). See [MAINTAINERS.md](./MAINTAINERS.md)
for instructions to keep up to date.

## v1.1.0

### Changes
* Bumped all firehose dependencies (from old versions)
* Logs are now more verbose by default output to STDERR
* Removed reader-node-merge-threshold-block-age (no more merging directly in the reader node, only the merger does it now)
* New 'ready' stat in prometheus for services

### UPGRADE Procedure

1. stop mindreader and delete its state file (note the previous block, your next start block will be the block preceding the 100-block boundary ex: 12345 -> startBlock: 12299)
2. stop the merger
3. delete remaining "one-block-files" from your oneblock folder
4. start new mindreader with flag changes and by specifying the new startBlock
  - rename all 'mindreader' mentions with just 'reader' (ex: `start reader-node` and `--mindreader-node-data-dir` becomes `--reader-node-data-dir`
  - rename `--common-one-blocks-store-url` to `--common-one-block-store-url` (no plural)
  - rename `--mindreader-node-debug-deep-mind` to `--reader-node-debug-firehose-logs`
  - remove `--mindreader-node-discard-after-stop-num`
  - remove `--mindreader-node-merge-threshold-block-age`
  - remove `--mindreader-node-blocks-chan-capacity`
5. start the new merger with flag changes (...) and watch it catch with live
  - rename `--common-one-blocks-store-url` to `--common-one-block-store-url` (no plural)
  - remove `--merger-time-between-store-lookups`
  - remove `--merger-writers-leeway`
  - remove `--merger-max-one-block-operations-batch-size`
  - add `--merger-time-between-store-pruning`
  - add `--merger-time-between-store-lookups`
6. replace the running relayer with the new one with flag changes
  - remove `--relayer-merger-addr`
  - remove `--relayer-buffer-size`
  - add `--common-one-block-store-url`
7. wait for it to become ready (takes some time... but it worked)
8. replace the firehose with the new one with flag changes:
  - remove `firehose-real-time-tolerance`
  - add `--common-one-block-store-url`

## v1.0.0

#### Flags and environment variables

* Renamed the `mindreader` application to `reader`
* Renamed `common-one-blocks-store-url` to `common-one-block-store-url`
* Renamed all the `mindreader-node-*` flags to `reader-node-*`





