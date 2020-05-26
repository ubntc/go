/*

This demo shows what happens when you write to the shared memory wihtout syncing.
You will read read bogus from memory!

*/

package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Integer stores integer values.
type Integer struct {
	i int32
	s string
}

// Store stores data.
type Store struct {
	s   string
	i   int32
	i64 int64
	s64 string
	o   Integer
	p   *Integer
}

// Err defines and race error
type Err int

// ProcFn modifies a value
type ProcFn func(*sync.WaitGroup)

// Race Error values
const (
	OK Err = iota
	Int32Error
	Int64Error
	Int32VarError
	Int64VarError
	String32Error
	String64Error
	StructError
	PointerError
	PanicError
)

var (
	prime    = int32(1777)
	prime64  = int64(prime)
	results  = make(chan Err, 1000)
	int32Var = randValue32()
	int64Var = int64(randValue32())
	strVar   = valueToString32(int32Var)
	str64Var = valueToString64(int64Var)
	store    = &Store{
		i: int32Var, s: strVar,
		i64: int64Var, s64: str64Var,
		o: Integer{int32Var, strVar},
		p: &Integer{int32Var, strVar},
	}
)

func check(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			results <- PanicError
		}
	}()
	defer wg.Done()
	errors := 0
	if int32Var%prime != 0 {
		results <- Int32VarError
		errors++
	}

	if int64Var%prime64 != 0 {
		results <- Int64VarError
		errors++
	}

	if store.i%prime != 0 {
		results <- Int32Error
		errors++
	}

	if store.i64%prime64 != 0 {
		results <- Int64Error
		errors++
	}

	if store.o.i%prime != 0 {
		results <- StructError
		errors++
	}

	if p := store.p; p == nil || p.i%prime != 0 {
		results <- PointerError
		errors++
	}

	if v, err := parse32(store.s); err != nil || v%prime != 0 {
		results <- String32Error
		errors++
	}

	if v, err := parse64(store.s64); err != nil || v%prime64 != 0 {
		results <- String64Error
		errors++
	}

	if v, err := parse32(store.o.s); err != nil || v%prime != 0 {
		results <- String32Error
		results <- StructError
		errors++
	}

	if v, err := parse32(store.p.s); err != nil || v%prime != 0 {
		results <- String32Error
		results <- PointerError
		errors++
	}

	if errors == 0 {
		results <- OK
	}
}

func randValue32() int32 {
	var v int32
	for v < prime {
		v = rand.Int31n(math.MaxInt32/prime) * prime
	}
	return v
}

func randValue64() int64 {
	var v int64
	for v < prime64 {
		v = rand.Int63n(math.MaxInt64/prime64) * prime64
	}
	return v
}

func parse32(v string) (int32, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func parse64(v string) (int64, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func valueToString32(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func valueToString64(v int64) string {
	return strconv.FormatInt(v, 10)
}

func modifyVars(wg *sync.WaitGroup) {
	defer wg.Done()
	int32Var = randValue32()
	int64Var = randValue64()
}

func modifyStore(wg *sync.WaitGroup) {
	defer wg.Done()
	v := randValue32()
	s := valueToString32(v)
	v64 := randValue64()
	s64 := valueToString64(v64)
	store.i = v
	store.s = s
	store.i64 = v64
	store.s64 = s64
	store.o = Integer{v, s}
	store.p = &Integer{v, s}
}

func proc(ctx context.Context) {
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			return
		default:
			wg.Add(3)
			go modifyVars(&wg)
			go modifyStore(&wg)
			go check(&wg)
			wg.Wait()
		}
	}
}

func countResults(ctx context.Context) {
	clr := "\r" + strings.Repeat(" ", 200) + "\r"
	var ok, int32Err, int64Err, int32VarErr, int64VarErr, str32Err, str64Err, structErr, ptrErr, panics uint64
	tick := time.Tick(10 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Printf("%sok: %d, intErr: %d|%d, intVarErr: %d|%d, strErr: %d|%d, strucErr: %d, ptrErr: %d, panics: %d",
				clr, ok, int32Err, int64Err, int32VarErr, int64VarErr, str32Err, str64Err, structErr, ptrErr, panics)
		case err := <-results:
			switch err {
			case OK:
				ok++
			case Int32Error:
				int32Err++
			case Int64Error:
				int64Err++
			case Int32VarError:
				int32VarErr++
			case Int64VarError:
				int64VarErr++
			case String32Error:
				str32Err++
			case String64Error:
				str64Err++
			case StructError:
				structErr++
			case PointerError:
				ptrErr++
			case PanicError:
				panics++
			}
			if ok+int64Err+str64Err > uint64(*maxResults) {
				fmt.Println()
				return
			}
		case <-ctx.Done():
			fmt.Println()
			return
		}
	}
}

var numProcs = flag.Int("n", 10, "num procs")
var maxResults = flag.Int("m", 1e6, "num procs")

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < *numProcs; i++ {
		go proc(ctx)
	}

	countResults(ctx)
	cancel()
	<-ctx.Done()
}
