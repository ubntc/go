package game

import (
	"time"

	"github.com/ubntc/go/games/gotris/game/geometry"
	"github.com/ubntc/go/games/gotris/game/tiles"
)

type GameState struct {
	Steps       uint
	Score       uint
	Speed       time.Duration
	Message     map[string]interface{}
	CurrentTile *tiles.Tile
	NextTile    *tiles.Tile
	Board       geometry.PointMap
	BoardPos    *geometry.Dim // BoardPos is the screen position of the top left corner of game screen.
}
