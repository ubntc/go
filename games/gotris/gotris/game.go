package gotris

import (
	"math"
	"time"
)

type Rules struct {
	// defines how many times the number of removed lines is multiplied with itself
	// to determine the score points for removed lines
	LineScoreExponent float64
	Width             int
	Height            int
	StepDuration      time.Duration
}

// Game stores the game state
type Game struct {
	Rules

	Steps       int
	CurrentTile *Tile
	NextTile    *Tile
	Tiles       []Tile

	// TODO: find better data structure to manage visual game state and occlusion processing
	Blocks [][]Point
}

func (g *Game) Advance() {
	g.Steps += 1
	if g.Move(g.CurrentTile, DirDown) {
		// tile moved down, game continues
		return
	}
	// tile stopped,
	g.ResolveTile(g.CurrentTile)
	g.Score(g.FindLines())
	g.CurrentTile = g.NextTile
	g.NextTile = NewTile(RandomTileType(), g.Width/2, g.Height)
}

func NewGame(rules Rules) *Game {
	return &Game{
		Rules: rules,
	}
}

// ResolveTile puts given tile's blocks into game grid.
func (g *Game) ResolveTile(t *Tile) {
}

// Move move a given tile one step in the given direction (U|R|D|L).
func (g *Game) Move(t *Tile, d Dir) bool {
	return true
}

// Rotate rotates (and if needed moves) the given tile in the given direction (L|R)
func (g *Game) Rotate(t *Tile, d Dir) bool {
	return true
}

// FindLines finds all completed lines and returns their row numbers.
func (g *Game) FindLines() []int {
	return nil
}

// Score coputes the score for a number of removed lines (given by line indexes).
func (g *Game) Score(lines []int) int {
	n := float64(len(lines))
	score := math.Pow(n, g.LineScoreExponent)
	return int(score)
}
