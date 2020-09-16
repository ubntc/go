package batbq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapacityRatios(t *testing.T) {
	testCases := []struct {
		name string
		cap  int
		sCap int
		dCap int
	}{
		{"big cap", 1000, 800, 200},
		{"med cap", 100, 80, 20},
		{"small cap", 10, 8, 2},
		{"zero cap", 0, 0, 0},
		{"1 cap", 1, 1, 0},
		{"2 cap", 2, 2, 0},
		{"5 cap", 5, 4, 1},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			stalledCap := StalledCapacity(c.cap)
			drainedCap := DrainedCapacity(c.cap)
			assert.Equal(t, c.sCap, stalledCap)
			assert.Equal(t, c.dCap, drainedCap)
		})
	}
}
