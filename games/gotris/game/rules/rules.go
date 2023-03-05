package rules

import (
	"time"

	"github.com/ubntc/go/games/gotris/common/geometry"
	"github.com/ubntc/go/games/gotris/common/platform"
)

const (
	SeedRandom platform.Seed = 0
)

var (
	TestRules = platform.Rules{
		BoardSize:   geometry.Dim{W: 5, H: 5},
		TickTime:    time.Millisecond,
		SpeedStep:   time.Nanosecond,
		MaxSteps:    10,
		PreviewSize: geometry.Dim{W: 3, H: 4},
		Seed:        123,
	}

	DefaultRules = platform.Rules{
		BoardSize:   geometry.Dim{W: 10, H: 20},
		TickTime:    time.Second,
		SpeedStep:   5 * time.Millisecond,
		MaxSteps:    0,
		PreviewSize: geometry.Dim{W: 5, H: 4},
	}
)
