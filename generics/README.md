This failed experiment wanted to create a generic `ForEach` method.
However you cannot cast. e.g., `[]int` to `[]interface{}` and need to use a `chan interface{}`
to pass sequential data around. The result is just to slow to provide any value.

ForEach results:

    BenchmarkForEach10-4      243576              5053 ns/op
    BenchmarkForEach1k-4        3158            373238 ns/op
    BenchmarkForEach1M-4           3         377289272 ns/op

Normal loop results:

    BenchmarkLoop10-4        1625445               727 ns/op
    BenchmarkLoop1k-4          19542             59471 ns/op
    BenchmarkLoop1M-4             19          61522294 ns/op
