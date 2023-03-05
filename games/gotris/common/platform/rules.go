package platform

import (
	"time"

	geom "github.com/ubntc/go/games/gotris/common/geometry"
)

type Seed int64

type Rules struct {
	TickTime    time.Duration // inital tick time to advance the game
	SpeedStep   time.Duration // how much to reduce the ticktime for anytime lines are scored
	MaxSteps    int           // max number of ticks the game can take (the default 0 is means infinity)
	BoardSize   geom.Dim      // size of the game board
	PreviewSize geom.Dim      // size of the preview box
	Seed        Seed          // Seed for randomization
}
