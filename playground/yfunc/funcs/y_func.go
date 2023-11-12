package funcs

// Func is a type representing a function that can compute a factorial
type Fn func(Fn, int64) int64

// Y is a function that takes a function f and returns a function that computes factorial
func YFn(f Fn) Fn {
	return func(g Fn, n int64) int64 {
		if n == 0 {
			return 1
		}
		return n * g(g, n-1)
		// return g(g, n-1) + g(g, n-2) // Fibonacci
	}
}

// Factorial is a function that computes the factorial of n using the Y combinator
func FactorialFunc(n int) any {
	f := YFn(nil)
	return f(f, int64(n))
}
