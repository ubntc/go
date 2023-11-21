package examples

import (
	"errors"
	"fmt"
	"log"
)

type step struct {
	Name string
	Func func() error
}

type Steps []step

func Step(name string, fn func() error) step {
	return step{name, fn}
}

func (steps *Steps) Run() error {
	if steps == nil {
		return nil
	}
	for i, s := range *steps {
		log.Printf("[Step %d] %s\n", i, s.Name)
		if err := s.Func(); err != nil {
			return errors.Join(err, fmt.Errorf("Step failed, step: %d, name: %s", i, s.Name))
		}
	}
	return nil
}

func (steps *Steps) Add(name string, fn func() error) *Steps {
	if steps == nil {
		*steps = make(Steps, 0)
	}
	*steps = append(*steps, step{name, fn})
	return steps
}
