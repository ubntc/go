package geometry

type Point struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

func OffsetPointsXY(points []Point, x, y int) []Point {
	res := make([]Point, len(points))
	for i, p := range points {
		res[i] = Point{p.X + x, p.Y + y}
	}
	return res
}

func OffsetPointsDir(points []Point, dir Dir) []Point {
	dx := 0
	dy := 0
	switch dir {
	case DirDown:
		dy = -1
	case DirUp:
		dy = +1
	case DirLeft:
		dx = -1
	case DirRight:
		dx = +1
	}
	return OffsetPointsXY(points, dx, dy)
}

func PointsInRange(points []Point, w, h int) bool {
	for _, p := range points {
		if p.X < 0 || p.X >= w {
			return false
		}
		if p.Y < 0 || p.Y >= h {
			return false
		}
	}
	return true
}
