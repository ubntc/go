package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ubntc/go/playground/yfunc/funcs"
)

const (
	MODE_STRUCT_64     = "struct64"
	MODE_STRUCT_BIGINT = "struct"
	MODE_ITER          = "iter"
	MODE_FUNC          = "func"
	MODE_LAMBDA        = "lambda"
)

var modes = []string{
	MODE_STRUCT_64,
	MODE_STRUCT_BIGINT,
	MODE_ITER,
	MODE_FUNC,
	MODE_LAMBDA,
}

func usage() {
	fmt.Println("arg1 must be one of:", modes)
}

func bitlen(v any) int64 {
	switch v := v.(type) {
	case int64, int:
		return 64
	case *big.Int:
		return int64(v.BitLen())
	}
	return -1
}

func main() {
	var fn func(int) any

	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case MODE_STRUCT_64:
		fn = funcs.StructFactorialInt64
	case MODE_STRUCT_BIGINT:
		fn = funcs.StructFactorial
	case MODE_ITER:
		fn = funcs.IterFactorial
	case MODE_FUNC:
		fn = funcs.FactorialFunc
	case MODE_LAMBDA:
		fn = funcs.FactorialCall
	default:
		usage()
		return
	}

	for _, n := range []int{5, 10, 20, 30, 40, 50, 100, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10} {
		result := fn(n)
		if w := bitlen(result); w <= 64 {
			fmt.Printf("%d: %d\n", n, result)
		} else {
			fmt.Printf("%d: too big (%d bits)\n", n, w)
		}
	}
}
