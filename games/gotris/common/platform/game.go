package platform

import (
	"time"

	geom "github.com/ubntc/go/games/gotris/common/geometry"
)

// Game is the common base struct that contains all data required
// for a `platform.Platform` to render and manage the Game.
type Game struct {
	Rules `json:"rules,omitempty"`

	Score       int           `json:"score,omitempty"`
	Speed       time.Duration `json:"speed,omitempty"`
	BoardPos    geom.Dim      `json:"board_pos,omitempty"`
	NextTile    *geom.Tile    `json:"next_tile,omitempty"`
	CurrentTile *geom.Tile    `json:"current_tile,omitempty"`
	Board       geom.PointMap `json:"board,omitempty"`

	Platform Platform `json:"-"`
}
