package main

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/klauspost/cpuid/v2"
)

type Store struct {
	value string
}

func (s *Store) Read() {
	v := s.value
	if _, err := strconv.ParseInt(v, 10, 64); err != nil {
		panic(err)
	}
}

func (s *Store) SetInt(v int64) {
	s.value = strconv.FormatInt(v, 10)
}

func main() {
	fmt.Println(CPU.BrandName)
	store := Store{value: "0"}
	for i := 0; i < 100000; i++ {
		go store.SetInt(time.Now().UnixNano())
		go store.Read()
	}
}
