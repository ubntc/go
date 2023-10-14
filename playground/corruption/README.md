# What could possibly "Go" wrong !?

This demo shows what can go wrong with Go shared memory if you do not properly
protect concurrent shared reads and writes.

The demo is part of the [ubntc/go](https://github.com/ubntc/go) Go playground
that collects and curates Go code examples for self-education and teaching
of good and bad practices.


Output on a RaspberryPi (arm):

    go run corruption.go
    ok: 900962, intErr: 58, strErr: 150, strucErr: 130, ptrErr: 28
    panic: runtime error: invalid memory address or nil pointer dereference


Output on a X1 Carbon (core i7):

    go run playground/corruption/corruption.go
    ok: 983638, intErr: 0, strErr: 227, strucErr: 79, ptrErr: 0


Output on an M1 Mac:

    go run playground/corruption/corruption.go
    ok: 989373, intErr: 0, strErr: 7562, strucErr: 7506, ptrErr: 4
    panic: runtime error: invalid memory address or nil pointer dereference


As you can see, in Go you can easily read bogus (unexpected values) and your program
may even crash if you read and write concurrently without proper `sync`.
