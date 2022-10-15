package game

import "time"

type Dim struct {
	Width  int
	Height int
}

type Seed int64

const (
	SeedRandom Seed = 0
)

type Rules struct {
	TickTime    time.Duration // inital tick time to advance the game
	SpeedStep   time.Duration // how much to reduce the ticktime for anytime lines are scored
	MaxSteps    int           // max number of ticks the game can take (the default 0 is means infinity)
	BoardSize   Dim           // size of the game board
	PreviewSize Dim           // size of the preview box
	Seed        Seed          // Seed for randomization
}

var (
	TestRules = Rules{
		BoardSize:   Dim{5, 5},
		TickTime:    time.Millisecond,
		SpeedStep:   time.Nanosecond,
		MaxSteps:    10,
		PreviewSize: Dim{3, 4},
		Seed:        123,
	}

	DefaultRules = Rules{
		BoardSize:   Dim{10, 20},
		TickTime:    time.Second,
		SpeedStep:   5 * time.Millisecond,
		MaxSteps:    0,
		PreviewSize: Dim{5, 4},
	}
)
