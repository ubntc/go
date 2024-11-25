# Parse Panic Demo
This demo shows how you should NOT program concurrent reads and writes.

Here is what happens when running this code on an **RPi4**:
```
$ go run main.go
panic: strconv.ParseInt: parsing "1669278767787263121": invalid syntax

goroutine 162111 [running]:
main.(*Store).Read(0x1d0c148)
        /home/pi/corruption/panictest/main.go:17 +0x68
created by main.main
        /home/pi/corruption/panictest/main.go:30 +0xac
exit status 2
```

Here is what you get on an **M1 Mac**:
```
$ go run parsepanic.go
Apple M1
panic: strconv.ParseInt: parsing "1669279266736419000": invalid syntax

goroutine 22971 [running]:
main.(*Store).Read(0x14000110290?)
        /Users/user/git/ubntc/go/playground/corruption/parsepanic/parsepanic.go:17 +0x58
created by main.main
        /Users/user/git/ubntc/go/playground/corruption/parsepanic/parsepanic.go:30 +0xa0
exit status 2
```

## Explanation

Looking at the code, you can see that the program never sets a non-Integer value:
```go
s.value = strconv.FormatInt(v, 10)
```
But `strconv.ParseInt` still observes an unparsable value. This is caused by the un-synced concurrent access to the string. Unexpected things do happen between the CPU, the CPU cache, and the RAM, when accessing data this way.

The panic message says `parsing "1669279266736419000"`, which is a valid value and seems to indicate that it should not have panicked. However, the message is not helpful, because the unparsable value was already overwritten at the time of logging it.

## Conclusion
* Always `sync` reads and writes!
* Use `go run -race` and `go test -race` as final gate keeper in your CI.
* Use linters that can uncover bad practices.
* Avoid (unprotected) global vars.
