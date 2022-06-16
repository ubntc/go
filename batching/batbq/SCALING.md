# Worker Behavior and Worker Scaling

------

**⚠️ Attention ⚠️**

Benchmarks are outdated and also a test against a non-batching naive streaming implementation is missing.

------


Batbq uses a blocking [`worker`](worker.go) to read data from the input channel into batches.
The worker will NOT stop reading data if the `Putter` or the message confirmation is stalled.
The calls to `Put(...)`, `Ack()`, and `Nack(error)` are fully asynchronous.
The sender must handle when to stop sending data; based on the number of unconfirmed messages.

Reading from the input channel into the current batch is done in one goroutine to avoid data races.
This can be a bottleneck and may require using more than one worker to read from the same channel.

If `BatcherConfig.AutoScale` is `true` the batcher will concurrently run more workers based on the
observed batch load and the [configured](config/config.go) `MinWorkers` and `MaxWorkers`.

Using multiple workers may result in more batches being collected and sent concurrently via
`output.Put(ctx, batch)`. However, all workers share the same `input <-chan batbq.Message` and the
same `output Putter`. Both, data source and output, must be concurrency-safe by supporting
concurrent calls of `Ack()`, `Nack(error)`, and `Put(ctx, batch)`.

## Benchmarks

You can play with the PubSub [publisher](_examples/publisher/main.go) and the
[ps2bq demo](_examples/ps2bq/main.go) to test which scaling parameters work best in your project.

1. Run `go run _examples/publisher/main.go` to start 100 concurrent pubsub writers that will push
   a total of 1000 [`Click`](_examples/ps2bq/clicks/clicks.go) events per second to a "click" topic
   in PubSub.

2. Concurrently run `time go run _examples/ps2bq/main.go -stats -c CAPACITY` with various `CAPACITY`
   values to see the effects of the chosen batch size on the worker scaling and CPU usage.


## Test Results: 1000 msg/s
Tests on a regular 4 core Linux laptop with a 100Mbit/s internet connection, using `MaxWorkers = 10`
and `MaxOutstandingMessages = 10000`, provided the following observations.

```
.----------------------------------------------------------------.
| input      | batch | workers  | output     | CPU usage         |
| rate       | size  |          | rate       | user   system cpu |
|----------------------------------------------------------------|
| 1000 msg/s | 10    | 10 (max) | 1000 msg/s | 9,34s  2,81s  31% |
| 1000 msg/s | 100   | 6        | 1000 msg/s | 5,63s  1,97s  18% |
| 1000 msg/s | 200   | 3        | 1000 msg/s | 5,33s  1,98s  17% |
| 1000 msg/s | 500   | 2        | 1000 msg/s | 4,72s  1,75s  16% |
| 1000 msg/s | 1000  | 1        | 1000 msg/s | 4,40s  1,53s  13% |
| 1000 msg/s | 2000  | 1        | 1000 msg/s | 4,32s  1,18s  13% |
| 1000 msg/s | 5000  | 1        | 1000 msg/s | 4,48s  1,45s  15% |
| 1000 msg/s | 10000 | 1        | 1000 msg/s | 4,58s  1,49s  15% |
`----------------------------------------------------------------´
```

As a result, using one worker with a batch capacity of 1000-2000 should be the preferred option for
an input rate of 1000 msg/s. Also, using much higher capacities does not have a big negative impact,
since the `DefaultFlushInterval` is `time.Second`.

## Test Results: 100 msg/s
```
.----------------------------------------------------------------.
| input      | batch | workers  | output     | CPU usage         |
| rate       | size  |          | rate       | user   system cpu |
|----------------------------------------------------------------|
| 100 msg/s  | 1000  | 1        | 100 msg/s  | 1,63s  0,38s  7%  |
| 100 msg/s  | 100   | 1        | 100 msg/s  | 1,57s  0,45s  7%  |
| 100 msg/s  | 10    | 6        | 100 msg/s  | 1,98s  0,54s  9%  |
`----------------------------------------------------------------´
```

## Test Results: 10 msg/s
```
.----------------------------------------------------------------.
| input      | batch | workers  | output     | CPU usage         |
| rate       | size  |          | rate       | user   system cpu |
|----------------------------------------------------------------|
| 10 msg/s   | 100   | 1        | 10 msg/s   | 1,57s  0,30s  6%  |
| 10 msg/s   | 20    | 1        | 10 msg/s   | 1,47s  0,33s  5%  |
| 10 msg/s   | 10    | 1        | 10 msg/s   | 1,43s  0,36s  6%  |
`----------------------------------------------------------------´
```
