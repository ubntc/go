package game

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	cmd "github.com/ubntc/go/games/gotris/game/controls"
	"github.com/ubntc/go/games/gotris/game/geometry"
	"github.com/ubntc/go/games/gotris/game/rules"
	"github.com/ubntc/go/games/gotris/game/scenes"
	"github.com/ubntc/go/games/gotris/game/tiles"
	"github.com/ubntc/go/games/gotris/input"
)

type GameConfig struct {
	rules.Rules

	GameOverScreenDuration time.Duration
}

var TestConfig = GameConfig{Rules: rules.TestRules}

// Game stores the game state
type Game struct {
	GameConfig
	GameState

	captureInput bool
	platform     Platform
	input        <-chan input.Input
}

func NewGame(cfg GameConfig, platform Platform) *Game {
	state := GameState{
		Board:    make(geometry.PointMap),
		Speed:    cfg.TickTime,
		BoardPos: geometry.NewDim(8, 0),
	}
	g := &Game{
		GameState:  state,
		GameConfig: cfg,

		captureInput: true,
		platform:     platform,
	}
	switch g.Seed {
	case rules.SeedRandom:
		rand.Seed(time.Now().Unix())
	default:
		rand.Seed(int64(g.Seed))
	}
	g.spawnTiles()
	return g
}

func (g *Game) DisableInput() { g.captureInput = false }

func (g *Game) Advance() (err error) {
	g.Steps += 1
	if err = g.move(g.CurrentTile, geometry.DirDown); err != nil {
		// tile hit another tile or the ground
		// split tile into blocks and check for lines
		// if the split fails the game is over (we are stuck somewhere on the top)
		if err := g.resolveTile(g.CurrentTile); err != nil {
			return err
		}
		lines := g.findLines()
		if len(lines) > 0 {
			g.Score += g.scoreLines(lines)
			g.removeLines(lines)
			g.Speed -= g.SpeedStep
		}
		g.spawnTiles()
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

func (g *Game) spawnTiles() {
	if g.NextTile == nil {
		g.NextTile = tiles.RandomTile()
	}
	// get next tile and create a new one
	t := g.NextTile
	g.NextTile = tiles.RandomTile()

	// move tile to the board
	dx := g.BoardSize.W / 2
	dy := g.BoardSize.H - 1
	t.Shift(dx, dy)
	g.CurrentTile = t
}

// resolveTile merged a given tile's blocks into the game blocks;
// returns and error if the tile cannot be merged (Game Over!).
func (g *Game) resolveTile(t *tiles.Tile) error {
	if g.Board.ContainsAny(t.Points()) {
		return errors.New("cannot resolve, tile collides with bocks")
	}
	tiles.MergeTile(t, g.Board)
	return nil
}

func (g *Game) updateCurrentTile(target tiles.Tile) error {
	if !geometry.PointsInRange(target.Points(), g.BoardSize.W, g.BoardSize.H+4) {
		return errors.New("tile not inside screen")
	}
	if g.Board.ContainsAny(target.Points()) {
		return errors.New("tile would collide with bocks")
	}
	g.CurrentTile.Update(target)
	return nil
}

// move move a given tile one step in the given direction (U|D|L|R).
func (g *Game) move(t *tiles.Tile, dir geometry.Dir) error {
	target := t.Moved(dir)
	if err := g.updateCurrentTile(target); err != nil {
		return errors.Wrap(err, "cannot move")
	}
	return nil
}

// rotate rotates (and if needed moves) the given tile in the given direction (CW|CCW)
func (g *Game) rotate(t *tiles.Tile, r geometry.Spin) error {
	rot := t.Rotated(r)
	if err := g.updateCurrentTile(rot); err != nil {
		return errors.Wrap(err, "cannot rotate")
	}
	return nil
}

func (g *Game) rotateAndMove(t *tiles.Tile, r geometry.Spin) (err error) {
	if err = g.rotate(t, r); err == nil {
		return err
	}
	if err = g.move(t, geometry.DirRight); err == nil {
		return g.rotate(t, r)
	}
	if err = g.move(t, geometry.DirLeft); err == nil {
		return g.rotate(t, r)
	}
	return err
}

func (g *Game) drop(t *tiles.Tile) error {
	for {
		if err := g.move(t, geometry.DirDown); err != nil {
			return err
		}
	}
}

// findLines finds all completed lines and returns their row numbers.
func (g *Game) findLines() (lines []int) {
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
func (g *Game) scoreLines(lines []int) uint {
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
	return uint(score)
}

func (g *Game) removeLines(lines []int) {
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

func (g *Game) runCommand(ctx context.Context, command cmd.Cmd, arg string) error {
	if dir := command.ToDir(); dir != geometry.DirUnkown {
		return g.move(g.CurrentTile, dir)
	}
	if spin := command.ToSpin(); spin != geometry.SpinUnknown {
		return g.rotateAndMove(g.CurrentTile, spin)
	}
	switch command {
	case cmd.Drop:
		_ = g.drop(g.CurrentTile)
		_ = g.Advance()
	case cmd.Help:
		g.showHelp(ctx)
	case cmd.Options:
		g.showOptions(ctx)
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
		g.platform.Options().Set(i - 1)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	return nil
}

func (g *Game) showScene(ctx context.Context, scene *scenes.Scene, timeout time.Duration) input.Input {
	g.platform.RenderScene(scene)
	return <-input.Await(ctx, g.input, timeout)
}
