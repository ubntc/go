## Basic Y Combinator in Go
and how to make it crash.

```bash➜  yfunc git:(master) ✗ go run yfunc.go func 2> >(head -n 3)
5: 120
10: 3628800
20: 2432902008176640000
30: -8764578968847253504
40: -70609262346240000
50: -3258495067890909184
100: 0
1000: 0
10000: 0
100000: 0
1000000: 0
10000000: 0
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp=0x14020160390 stack=[0x14020160000, 0x14040160000]
fatal error: stack overflow
```
