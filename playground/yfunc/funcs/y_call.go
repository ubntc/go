package funcs

// Unified function interface
type Func interface {
	Func(Func, int) int
}

// FuncImpl to convert a function to the Func interface
type FuncImpl func(Func, int) int

func (f FuncImpl) Func(g Func, n int) int {
	return f(g, n)
}

// Y combinator function
func Y(f Func) Func {
	return FuncImpl(func(g Func, n int) int {
		return f.Func(g, n)
	})
}

// Factorial function implementation
type FactorialImplementer struct{}

func (f FactorialImplementer) Func(g Func, n int) int {
	if n == 0 {
		return 1
	}
	return n * g.Func(g, n-1)
}

func FactorialCall(n int) any {
	YFactorial := Y(FactorialImplementer{})
	result := YFactorial.Func(YFactorial, n)
	return result
}
