# Ubntc Go Projects

This repository hosts the following projects.

## [Go-cli: ubntc/go/cli](/cli)
Basic CLI-enhancements for your Go-services, incl. input commands, human-friendly logging, and
OS-signal handling.

## [BatSub: ubntc/go/batsub](/batsub)
Capacity and interval-based batch reading of pubsub messages.

## [Go-scripts: ubntc/go/scripts](/scripts)
Reusable build scrips and utils for managing Go code in this monorepo.

## [Gophers: ubntc/go/gophers](/gophers)
Gophers art and vector graphics.

## [Playground: ubntc/go/playground](/playground)
Experiments and code for learning and understanding the pitfalls of Go, esp. regarding concurrency.

## Monorepo approach
The Go code in this monorepo is managed in subpackages, such as [ubntc/cli](cli).
These subpackages have their own `go.mod` file, which makes `go mod` exclude their dependencies
from the root [go.mod]().

Maturing packages are given their own `go.mod`, `README.md`, and a copy of the `LICENSE` as soon as
they are tested thoroughly and fulfill a specific purpose on their own.

All other Go code is to be considered highly experimental and is owned by the root project from
where it could be vendored for testing purposes only.