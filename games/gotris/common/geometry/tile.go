package geometry

type Typ string

type Tile struct {
	Typ         Typ     `json:"type,omitempty"`
	Orientation Dir     `json:"orientation,omitempty"`
	Points      []Point `json:"points,omitempty"`
	Center      int     `json:"center,omitempty"`

	// Drawings map  for looking up tile layouts
	Drawings Drawings `json:"drawings,omitempty"`
}

// Drawings provides a data structure to store drawing variants of each tile.
type Drawings map[Typ][]string

func NewTile(typ Typ, x, y int, drawings Drawings) *Tile {
	points, ori, center := drawings.PointsForType(typ, DirUp)
	return &Tile{
		Typ:         typ,
		Orientation: ori,
		Points:      OffsetPointsXY(points, x, y),
		Center:      center,
		Drawings:    drawings,
	}
}

func (t *Tile) SetPoints(points []Point, ori Dir, center int) {
	t.Points = points
	t.Orientation = ori
	t.Center = center
}

func (t *Tile) Move(dx, dy int) {
	t.Points = OffsetPointsXY(t.Points, dx, dy)
}

func (t *Tile) CenterPoint() int {
	return t.Center
}

func (t *Tile) Orientations() []string {
	return t.Drawings[t.Typ]
}

func (t *Tile) RotatedPoints(spin Spin) (points []Point, orientation Dir, center int) {
	numOris := Dir(len(t.Orientations()))
	ori := Dir(t.Orientation)
	switch spin {
	case SpinLeft:
		ori = (ori + numOris - 1) % numOris
	case SpinRight:
		ori = (ori + 1) % numOris
	}

	// obtain zero-centered points using the new orientation
	points, dir, center := t.Drawings.PointsForType(t.Typ, ori)

	// move the points in place
	x, y := t.Position()
	points = OffsetPointsXY(points, x, y)

	return points, dir, center
}

func (t *Tile) Position() (x int, y int) {
	return t.Points[t.Center].X, t.Points[t.Center].Y
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
	blocks.SetAll(t.Points, string(t.Typ))
}
