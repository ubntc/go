package rules

import (
	"time"

	"github.com/ubntc/go/games/gotris/game/geometry"
)

type Seed int64

const (
	SeedRandom Seed = 0
)

type Rules struct {
	TickTime    time.Duration // inital tick time to advance the game
	SpeedStep   time.Duration // how much to reduce the ticktime for anytime lines are scored
	MaxSteps    uint          // max number of ticks the game can take (the default 0 is means infinity)
	BoardSize   geometry.Dim  // size of the game board
	PreviewSize geometry.Dim  // size of the preview box
	Seed        Seed          // Seed for randomization
}

var (
	TestRules = Rules{
		BoardSize:   geometry.Dim{W: 5, H: 5},
		TickTime:    time.Millisecond,
		SpeedStep:   time.Nanosecond,
		MaxSteps:    10,
		PreviewSize: geometry.Dim{W: 3, H: 4},
		Seed:        123,
	}

	DefaultRules = Rules{
		BoardSize:   geometry.Dim{W: 10, H: 20},
		TickTime:    time.Second,
		SpeedStep:   5 * time.Millisecond,
		MaxSteps:    0,
		PreviewSize: geometry.Dim{W: 5, H: 4},
	}
)
