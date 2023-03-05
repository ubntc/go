package platform

import (
	"time"

	geom "github.com/ubntc/go/games/gotris/common/geometry"
)

// Game is the common base struct that contains all data required
// for a `platform.Platform` to render and manage the Game.
type Game struct {
	Rules

	Score       int
	Speed       time.Duration
	BoardPos    geom.Dim
	NextTile    *geom.Tile
	CurrentTile *geom.Tile
	Board       geom.PointMap

	Platform Platform
}
