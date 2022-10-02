package game

type Point struct {
	X int
	Y int
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

func TopLeft(points []Point) Point {
	x := points[0].X
	y := points[0].Y
	for _, p := range points {
		if p.X < x {
			x = p.X
		}
		if p.Y < y {
			y = p.Y
		}
	}
	return Point{x, y}
}

func Add(points []Point, x, y int) []Point {
	res := make([]Point, len(points))
	for i, p := range points {
		res[i] = Point{p.X + x, p.Y + y}
	}
	return res
}

func Pivot(points []Point) []Point {
	res := make([]Point, len(points))
	for i, p := range points {
		res[i] = Point{p.Y, p.X}
	}
	return res
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
