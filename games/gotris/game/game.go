package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

// Game stores the game state
type Game struct {
	Rules

	Steps   int
	Score   int
	Speed   time.Duration
	Message map[string]interface{}

	CurrentTile *Tile
	NextTile    *Tile

	Board PointMap
}

func NewGame(rules Rules) *Game {
	g := &Game{
		Rules: rules,
		Board: make(PointMap),
		Speed: rules.TickTime,
	}
	switch g.Seed {
	case SeedRandom:
		rand.Seed(time.Now().Unix())
	default:
		rand.Seed(int64(g.Seed))
	}
	g.SpawnTiles()
	return g
}

func (g *Game) SpawnTiles() {
	if g.NextTile == nil {
		g.NextTile = RandomTile()
	}
	// get next tile and create a new one
	t := g.NextTile
	g.NextTile = RandomTile()

	// move tile to the board
	dx := g.BoardSize.Width / 2
	dy := g.BoardSize.Height - 1
	t.points = OffsetPointsXY(t.points, dx, dy)
	g.CurrentTile = t
}

func (g *Game) Advance() (err error) {
	g.Steps += 1
	if err = g.Move(g.CurrentTile, DirDown); err != nil {
		// tile hit another tile or the ground
		// split tile into blocks and check for lines
		// if the split fails the game is over (we are stuck somewhere on the top)
		if err := g.ResolveTile(g.CurrentTile); err != nil {
			return err
		}
		lines := g.FindLines()
		if len(lines) > 0 {
			g.Score += g.ScoreLines(lines)
			g.RemoveLines(lines)
			g.Speed -= g.SpeedStep
		}
		g.SpawnTiles()
	}
	return nil
}

func (g *Game) AdvanceBy(steps int) error {
	for {
		if steps <= 0 {
			return nil
		}
		steps -= 1
		if err := g.Advance(); err != nil {
			return err
		}
	}
}

// ResolveTile merged a given tile's blocks into the game blocks;
// returns and error if the tile cannot be merged (Game Over!).
func (g *Game) ResolveTile(t *Tile) error {
	if g.Board.ContainsAny(t.Points()) {
		return errors.New("cannot resolve, tile collides with bocks")
	}
	MergeTile(t, g.Board)
	return nil
}

func (g *Game) ModifyTile(t *Tile, points []Point, ori Dir, center int) error {
	if !PointsInRange(points, g.BoardSize.Width, g.BoardSize.Height+4) {
		return errors.New("tile not inside screen")
	}
	if g.Board.ContainsAny(points) {
		return errors.New("tile would collide with bocks")
	}
	t.SetPoints(points, ori, center)
	return nil
}

// Move move a given tile one step in the given direction (U|D|L|R).
func (g *Game) Move(t *Tile, dir Dir) error {
	points := OffsetPointsDir(t.Points(), dir)
	if err := g.ModifyTile(t, points, t.orientation, t.center); err != nil {
		return errors.Wrap(err, "cannot move")
	}
	return nil
}

// Rotate rotates (and if needed moves) the given tile in the given direction (CW|CCW)
func (g *Game) Rotate(t *Tile, r Spin) error {
	points, ori, center := t.RotatedPoints(r)
	if err := g.ModifyTile(t, points, ori, center); err != nil {
		return errors.Wrap(err, "cannot rotate")
	}
	return nil
}

func (g *Game) RotateAndMove(t *Tile, r Spin) (err error) {
	if err = g.Rotate(t, r); err == nil {
		return err
	}
	if err = g.Move(t, DirRight); err == nil {
		return g.Rotate(t, r)
	}
	if err = g.Move(t, DirLeft); err == nil {
		return g.Rotate(t, r)
	}
	return err
}

func (g *Game) Drop(t *Tile) error {
	for {
		if err := g.Move(t, DirDown); err != nil {
			return err
		}
	}
}

// FindLines finds all completed lines and returns their row numbers.
func (g *Game) FindLines() (lines []int) {
	rows := make(map[int]int)
	for p := range g.Board {
		rows[p.Y] += 1
		if rows[p.Y] == g.BoardSize.Width {
			lines = append(lines, p.Y)
		}
	}
	return
}

// Score coputes the score for a number of removed lines (given by line indexes).
func (g *Game) ScoreLines(lines []int) int {
	n := len(lines)
	lineFactor := 10
	// Speed Bonus Points per Line
	// (1000 ms - 995 ms) / 5 ms = 5 / 5 = 1
	// (1000 ms - 970 ms) / 5 ms = 30 / 5 = 6
	// (1000 ms - 900 ms) / 5 ms = 100 / 5 = 25
	// (1000 ms - 500 ms) / 5 ms = 500 / 5 = 100
	// (1000 ms - 100 ms) / 5 ms = 900 / 5 = 180
	// (1000 ms -  20 ms) / 5 ms = 980 / 5 = 196
	speedFactor := int((g.TickTime-g.Speed)/g.SpeedStep) + 1
	score := n * n * lineFactor * speedFactor * 10
	return int(score)
}

func (g *Game) RemoveLines(lines []int) {
	points := g.Board.PointsList(g.BoardSize.Width, g.BoardSize.Height)
	for _, y := range lines {
		fmt.Println("remove", y)
		// mark line to be deleted
		points[y] = nil
	}
	reduced := make([][]string, 0, len(points)-len(lines))
	for _, l := range points {
		if l != nil {
			reduced = append(reduced, l)
		}
	}
	g.Board.Clear()
	g.Board.SetPoints(reduced)
}

func (g *Game) Dump() {
	fmt.Println("Blocks", g.Board)
	fmt.Println("CurrentTile", g.CurrentTile)
	fmt.Println("NextTile", g.NextTile)
}

func (g *Game) RunCommand(cmd Cmd) error {
	if dir := cmd.ToDir(); dir != DirUnkown {
		return g.Move(g.CurrentTile, dir)
	}
	if spin := cmd.ToSpin(); spin != SpinUnknown {
		return g.RotateAndMove(g.CurrentTile, spin)
	}
	if cmd == CmdDrop {
		g.Drop(g.CurrentTile)
		g.Advance()
	}
	return fmt.Errorf("unknown command: %s", cmd)
}
