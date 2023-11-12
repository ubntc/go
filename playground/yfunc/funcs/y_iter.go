package funcs

import "math/big"

type Fii func(n int) *big.Int

func IterFactorial(n int) any {
	iter := func(f func(int) *big.Int, n int) *big.Int {
		result := big.NewInt(1)
		for ; n > 0; n-- {
			_ = result.Mul(result, big.NewInt(int64(n)))
		}
		return result
	}
	// Wrap iter in a function that matches the expected signature
	var wrapper Fii
	wrapper = func(n int) *big.Int {
		return iter(wrapper, n)
	}
	return wrapper(n)
}
