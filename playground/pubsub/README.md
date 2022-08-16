# Basic Pub/Sub
Implementation of a naive pub/sub system in Go.

## Features
* publish messages
* subscribe to messages (no topics)
* backfill old messages for new subscribers
* basic message retention strategy
* concurrency safe (tested with `--race`)
