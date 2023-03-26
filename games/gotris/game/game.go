package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	geom "github.com/ubntc/go/games/gotris/common/geometry"
	"github.com/ubntc/go/games/gotris/common/input"
	"github.com/ubntc/go/games/gotris/common/platform"
	cmd "github.com/ubntc/go/games/gotris/game/controls"
	"github.com/ubntc/go/games/gotris/game/rules"
)

// Game stores the game state
type Game struct {
	platform.Game

	Steps int

	GameOverScreenDuration time.Duration

	Message map[string]interface{}

	CaptureInput bool

	input <-chan *input.Input
}

func NewGame(gameRules platform.Rules, renderngPlatform platform.Platform) *Game {
	g := &Game{
		Game: platform.Game{
			Rules:    gameRules,
			BoardPos: *geom.NewDim(8, 0),
			Platform: renderngPlatform,
		},
		CaptureInput: true,
	}
	g.Init()
	return g
}

// Init resets the game board and speed and spawns the first tile.
func (g *Game) Init() {
	g.Board = make(geom.PointMap)
	g.Speed = g.TickTime
	g.Score = 0

	switch g.Seed {
	case rules.SeedRandom:
		rand.Seed(time.Now().Unix())
	default:
		rand.Seed(int64(g.Seed))
	}
	g.SpawnTile()
}

// SpawnTile moves the next tile in the preview onto the game board
// and creates a tile in the preview.
func (g *Game) SpawnTile() {
	if g.NextTile == nil {
		g.NextTile = RandomTile()
	}
	// get next tile and create a new one
	t := g.NextTile
	g.NextTile = RandomTile()

	// move tile to the board
	dx := g.BoardSize.W / 2
	dy := g.BoardSize.H - 1
	t.Move(dx, dy)
	g.CurrentTile = t
}

func (g *Game) Advance() (err error) {
	g.Steps += 1
	if err = g.Move(g.CurrentTile, geom.DirDown); err != nil {
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
		g.SpawnTile()
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
func (g *Game) ResolveTile(t *geom.Tile) error {
	if g.Board.ContainsAny(t.Points) {
		return errors.New("cannot resolve, tile collides with bocks")
	}
	geom.MergeTile(t, g.Board)
	return nil
}

func (g *Game) ModifyTile(t *geom.Tile, points []geom.Point, ori geom.Dir, center int) error {
	if !geom.PointsInRange(points, g.BoardSize.W, g.BoardSize.H+4) {
		return errors.New("tile not inside screen")
	}
	if g.Board.ContainsAny(points) {
		return errors.New("tile would collide with bocks")
	}
	t.SetPoints(points, ori, center)
	return nil
}

// Move move a given tile one step in the given direction (U|D|L|R).
func (g *Game) Move(t *geom.Tile, dir geom.Dir) error {
	points := geom.OffsetPointsDir(t.Points, dir)
	if err := g.ModifyTile(t, points, t.Orientation, t.CenterPoint()); err != nil {
		return errors.Wrap(err, "cannot move")
	}
	return nil
}

// Rotate rotates (and if needed moves) the given tile in the given direction (CW|CCW)
func (g *Game) Rotate(t *geom.Tile, r geom.Spin) error {
	points, ori, center := t.RotatedPoints(r)
	if err := g.ModifyTile(t, points, ori, center); err != nil {
		return errors.Wrap(err, "cannot rotate")
	}
	return nil
}

func (g *Game) RotateAndMove(t *geom.Tile, r geom.Spin) (err error) {
	if err = g.Rotate(t, r); err == nil {
		return err
	}
	if err = g.Move(t, geom.DirRight); err == nil {
		return g.Rotate(t, r)
	}
	if err = g.Move(t, geom.DirLeft); err == nil {
		return g.Rotate(t, r)
	}
	return err
}

func (g *Game) Drop(t *geom.Tile) error {
	for {
		if err := g.Move(t, geom.DirDown); err != nil {
			return err
		}
	}
}

// FindLines finds all completed lines and returns their row numbers.
func (g *Game) FindLines() (lines []int) {
	rows := make(map[int]int)
	for p := range g.Board {
		rows[p.Y] += 1
		if rows[p.Y] == g.BoardSize.W {
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
	points := g.Board.PointsList(g.BoardSize.W, g.BoardSize.H)
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

func (g *Game) RunCommand(command cmd.Cmd, arg string) error {
	if dir := command.ToDir(); dir != geom.DirUnkown {
		return g.Move(g.CurrentTile, dir)
	}
	if spin := command.ToSpin(); spin != geom.SpinUnknown {
		return g.RotateAndMove(g.CurrentTile, spin)
	}
	switch command {
	case cmd.Drop:
		g.Drop(g.CurrentTile)
		g.Advance()
	case cmd.Help:
		g.showHelp()
	case cmd.Options:
		g.showOptions()
	case cmd.MoveBoardLeft:
		if g.BoardPos.W > 1 {
			g.BoardPos.W -= 1
		}
	case cmd.MoveBoardRight:
		g.BoardPos.W += 1
	case cmd.MoveBoardUp:
		if g.BoardPos.H > 1 {
			g.BoardPos.H -= 1
		}
	case cmd.MoveBoardDown:
		g.BoardPos.H += 1
	case cmd.SelectMode:
		i, _ := strconv.Atoi(arg)
		g.Platform.Options().Set(i - 1)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	return nil
}

func (g *Game) ShowScene(scene platform.Scene, timeout time.Duration) *input.Input {
	g.Platform.RenderScene(scene)
	return input.AwaitInput(g.input, timeout)
}
