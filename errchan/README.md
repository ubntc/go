## errchan: idomatic, fast, and memory-efficient Go for concurrent multi-error handling.

# Benchmark Results

The benchmarks in this submodule compare `errchan` vs. `multierror`, `multierr`, and `errgroup`
in a simple concurrent usage scenario with an error rate of 10%.

```
BenchmarkMultiError10-4       	  181548	      6015 ns/op	     384 B/op	      10 allocs/op
BenchmarkMultiError1000-4     	    6004	    210554 ns/op	    6033 B/op	     116 allocs/op
BenchmarkMultiError100000-4   	      57	  22935904 ns/op	 1628013 B/op	   11642 allocs/op
BenchmarkMultiErr10-4         	  273561	      3868 ns/op	     320 B/op	       8 allocs/op
BenchmarkMultiErr1000-4       	    5726	    214545 ns/op	   10912 B/op	     213 allocs/op
BenchmarkMultiErr100000-4     	      55	  22082171 ns/op	 1472408 B/op	   20049 allocs/op
BenchmarkErrorGroup10-4       	  259068	      4373 ns/op	     456 B/op	      15 allocs/op
BenchmarkErrorGroup1000-4     	    4765	    286162 ns/op	   17880 B/op	    1104 allocs/op
BenchmarkErrorGroup100000-4   	      39	  28477494 ns/op	 1760281 B/op	  110004 allocs/op
BenchmarkErrChan10-4          	  216470	      4713 ns/op	     552 B/op	       8 allocs/op
BenchmarkErrChan1000-4        	    5157	    236809 ns/op	   22424 B/op	     114 allocs/op
BenchmarkErrChan100000-4      	      51	  23074991 ns/op	 2591960 B/op	   10026 allocs/op
```

# Observations
All approaches have their performance in the same order of magnitude.

1. `multierror` has the lowest memory consumption
   but requires a `sync.Mutex` to ensure concurrency safety.
2. `multierr` has 2nd lowest memory consumption
   but also requires a `sync.Mutex` to ensure concurrency safety.
3. `errgroup` is slowest, has 10x more allocs/op than the others,
   is concurrency-safe, but does not capture all errors.
4. `errchan` is 20% faster than `errgroup`, is concurrency-safe, captures all errors,
   has the lowest allocs/op, but has 50% higher memory consumption.

# Conclusions
For concurrent workloads where `errgroup` is useful, `errchan` can provide some speed benefit at the cost of extra memory.
It can also safely capture all errors, which `errgroup` can't.

Moreover, `errchan` can be used with any `go func()` instead of `grp.Go(func() error)`,
which allows `loopclosure` to detect misused loop variables and avoid programming errors.
However, it requires proper `sync.WaitGroup` usage and a predefined size that is bigger
than the number of expected errors.