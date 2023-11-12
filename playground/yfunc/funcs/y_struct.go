package funcs

import "math/big"

type recursiveBigInt struct {
	f func(recursiveBigInt, int) *big.Int
}

type recursiveInt struct {
	f func(recursiveInt, int) int64
}

func StructFactorial(n int) any {
	y := recursiveBigInt{
		f: func(y recursiveBigInt, n int) *big.Int {
			if n == 0 {
				return big.NewInt(1)
			}
			return big.NewInt(0).Mul(big.NewInt(int64(n)), y.f(y, n-1))
		},
	}
	return y.f(y, n)
}

func StructFactorialInt64(n int) any {
	y := recursiveInt{
		f: func(y recursiveInt, n int) int64 {
			if n == 0 {
				return 1
			}
			return int64(n) * y.f(y, n-1)
		},
	}
	return y.f(y, n)
}
