package funcs

// Define function types
type (
	// FnInt64 is a specific type of recursive function that implements some data logic
	// and decides when to recursively call itself. However, it does not have a direct
	// reference to itself but must use a function pointer from a closure or variable.
	FnInt64 func(int64) int64

	// YFuncGenerator is a function that takes a FnInt64 function returns a function to call
	// the given FnInt64 recursively.
	YFuncGenerator func(FnInt64) FnInt64
)

// YRecurse create a new function for every iteration and then calls this function.
func YRecurse(newFunc YFuncGenerator) FnInt64 {
	recurse := func(iteration int64) int64 {
		return YRecurse(newFunc)(iteration)
	}

	return newFunc(recurse)
}

func newFactorialFunc(f FnInt64) FnInt64 {
	return func(n int64) int64 {
		if n == 0 {
			return 1
		}
		return n * f(n-1)
	}
}

var FactorialSelfApply = YRecurse(newFactorialFunc)
