package geometry

type Typ string

type Tile struct {
	typ         Typ
	orientation Dir
	points      []Point
	center      int

	// drawings map  for looking up tile layouts
	drawings Drawings
}

// Drawings provides a data structure to store drawing variants of each tile.
type Drawings map[Typ][]string

func NewTile(typ Typ, x, y int, drawings Drawings) *Tile {
	points, ori, center := drawings.PointsForType(typ, DirUp)
	return &Tile{
		typ:         typ,
		orientation: ori,
		points:      OffsetPointsXY(points, x, y),
		center:      center,
		drawings:    drawings,
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

func (t *Tile) Move(dx, dy int) {
	t.points = OffsetPointsXY(t.points, dx, dy)
}

func (t *Tile) CenterPoint() int {
	return t.center
}

func (t *Tile) Orientation() Dir {
	return t.orientation
}

func (t *Tile) Orientations() []string {
	return t.drawings[t.typ]
}

func (t *Tile) RotatedPoints(spin Spin) (points []Point, orientation Dir, center int) {
	numOris := Dir(len(t.Orientations()))
	ori := Dir(t.orientation)
	switch spin {
	case SpinLeft:
		ori = (ori + numOris - 1) % numOris
	case SpinRight:
		ori = (ori + 1) % numOris
	}

	// obtain zero-centered points using the new orientation
	points, dir, center := t.drawings.PointsForType(t.typ, ori)

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

// PointsForType returns the points for a given Typ and orientation.
//
// A tiles points are centered around a center point, determined by the drawing
// instructions of the tile type and the given orientation. The resulting orientation
// and the index of the center point are returned as additional return values.
// All three values are needed to fully specify a tile on the game board.
func (d Drawings) PointsForType(typ Typ, orientation Dir) ([]Point, Dir, int) {
	points := make([]Point, 0, 4)
	x := 0
	y := 0
	orientation = orientation % Dir(len(d))
	center := 0
	for _, instruction := range d[typ][orientation] {
		switch instruction {
		case 'c':
			center = len(points)
			fallthrough
		case 'x':
			points = append(points, *NewPoint(x, y))
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

func MergeTile(t *Tile, blocks PointMap) {
	blocks.SetAll(t.Points(), string(t.Typ()))
}
