package game

import "math/rand"

type Typ string

const (
	TypB Typ = "B"
	TypI Typ = "I"
	TypL Typ = "L"
	TypJ Typ = "J"
	TypT Typ = "T"
	TypS Typ = "S"
	TypZ Typ = "Z"
)

var Tiles = []struct {
	Typ Typ
	// Block defines a single block of a tile
	Block string
	// Drawings defines how the tile is drawn using a small DSL:
	//  'x' = draw + move right
	//  'c' = draw + move right + use as center point
	//  ' ' = move right
	//  ',' = next row
	Drawings []string
}{
	// All colors, orientations, and center blocks for 7 lines of code! 🤯
	{TypB, "🟫", []string{"xc,xx", "xc,xx", "xc,xx", "xc,xx"}},
	{TypI, "🟨", []string{"x,c,x,x", "xcxx", "x,x,c,x", "xxcx"}},
	{TypL, "🟥", []string{"x,c,xx", "xcx,x", "xx, c, x", "  x,xcx"}},
	{TypJ, "🟧", []string{" x, c,xx", "x,xcx", "xx,c,x", "xcx,  x"}},
	{TypT, "🟩", []string{"xcx, x", " x,xc, x", " x,xcx", "x,cx,x"}},
	{TypS, "🟦", []string{" cx,xx", "x,xc, x", " xx,xc", "x,cx, x"}},
	{TypZ, "🟪", []string{"xc, xx", " x,xc,x", "xx, cx", " x,cx,x"}},
}

//  More coloring ideas:
//  ░▒▓█  ▤▥▦▧▨▩▣
//  🟥🟦🟧🟨🟩🟪🟫

var (
	// Drawings by type
	Drawings = make(map[Typ][]string, len(Tiles))
	// Blocks by type
	Blocks = make(map[Typ]string, len(Tiles))
)

func init() {
	for _, spec := range Tiles {
		Drawings[spec.Typ] = spec.Drawings
		Blocks[spec.Typ] = spec.Block
	}
}

type Tile struct {
	typ         Typ
	orientation Dir
	points      []Point
	center      int
}

func NewTile(typ Typ, x, y int) *Tile {
	points, ori, center := PointsForType(typ, DirUp)
	return &Tile{
		typ:         typ,
		orientation: ori,
		points:      OffsetPointsXY(points, x, y),
		center:      center,
	}
}

func (t *Tile) Points() []Point {
	return t.points
}

func (t *Tile) SetPoints(points []Point, ori Dir, center int) {
	t.points = points
	t.orientation = ori
	t.center = center
}

func (t *Tile) Orientations() []string {
	return Drawings[t.typ]
}

func (t *Tile) RotatedPoints(spin Spin) (points []Point, orientation Dir, center int) {
	oris := t.Orientations()
	numOris := Dir(len(oris))
	ori := Dir(t.orientation)
	switch spin {
	case SpinLeft:
		ori = (ori + numOris - 1) % numOris
	case SpinRight:
		ori = (ori + 1) % numOris
	}

	// obtain zero-centered points using the new orientation
	points, dir, center := PointsForType(t.typ, ori)

	// move the points in place
	x, y := t.Position()
	points = OffsetPointsXY(points, x, y)

	return points, dir, center
}

func (t *Tile) Center() Point {
	return t.points[t.center]
}

func (t *Tile) Position() (x int, y int) {
	return t.Center().X, t.Center().Y
}

func (t *Tile) Typ() Typ {
	return t.typ
}

func RandomTileType() Typ {
	i := rand.Int() % len(Tiles)
	return Tiles[i].Typ
}

func RandomTile() *Tile {
	return NewTile(RandomTileType(), 0, 0)
}

// PointsForType returns the points for a given Typ and orientation.
//
// A tiles points are centered around a center point, determined by the drawing
// instructions of the tile type and the given orientation. The resulting orientation
// and the index of the center point are returned as additional return values.
// All three values are needed to fully specify a tile on the game board.
func PointsForType(typ Typ, orientation Dir) ([]Point, Dir, int) {
	points := make([]Point, 0, 4)
	x := 0
	y := 0
	orientation = orientation % Dir(len(Drawings[typ]))
	center := 0
	for _, instruction := range Drawings[typ][orientation] {
		switch instruction {
		case 'c':
			center = len(points)
			fallthrough
		case 'x':
			points = append(points, Point{x, y})
			fallthrough
		case ' ':
			x += 1
		case ',':
			x = 0
			y -= 1
		}
	}
	c := points[center]
	return OffsetPointsXY(points, -c.X, -c.Y), Dir(orientation), center
}
