# From Racy to Non-Racy

Run the `bad` example first.
```
go run --race playground/bogus/bogus.go -bogus
==================
WARNING: DATA RACE
Read at 0x00c0000a0220 by goroutine 17: bogus.go:37 +0x85
Previous write at 0x00c0000a0220 by goroutine 16: bogus.go:25 +0xb6
...
==================
read: name = 9
read: name = 3
read: name = 3
read: name = 5
read: name = 6
read: name = 6
read: name = 4
read: name = 7
read: name = 7
Found 1 data race(s)
exit status 66
```
Race conditions are a hint that concurrent reads and writes are not proteced with `sync` primitives.


Now run the `good` example:
```
go run --race playground/bogus/bogus.go
read: name = 9
read: name = 3
read: name = 3
read: name = 5
read: name = 6
read: name = 6
read: name = 4
read: name = 7
read: name = 7
read: name = 2
```
All good!