[![Go Report Card](https://goreportcard.com/badge/github.com/ubntc/go/batching/batsub)](https://goreportcard.com/report/github.com/ubntc/go/batching/batsub)
[![cover-badge](https://img.shields.io/badge/coverage-89%25-brightgreen.svg?longCache=true&style=flat)](Makefile#8)

# Batched PubSub Reader
This package implements batched reading of messages from a `pubsub.Subscription`.
It provides a `BatchedSubscription` with a `ReceiveBatch` method to read messages in batches
based on a given batch capacity or batching interval.

## Usage

```golang
capacity := batsub.Capacity(1000)
interval := batsub.FlushInterval(time.Second)

sub := batsub.NewBatchedSubscription(subscription, capacity, interval)
err := sub.ReceiveBatch(ctx, func(ctx context.Context, messages []*pubsub.Message){
    // handle batch of messages using a batch-processing library
    errors := mylib.BatchProcessMessages(messages)
    for i, err := errors {
        if err != nil {
            // TODO: handle error
            continue
        }
        messages[i].Ack()
    }
})
```
