package game

type PointMap map[Point]string

func (m PointMap) ContainsAny(points []Point) bool {
	for _, p := range points {
		if _, ok := m[p]; ok {
			return true
		}
	}
	return false
}

func (m PointMap) Collides(points []Point, dir Dir) bool {
	points = OffsetPointsDir(points, dir)
	return m.ContainsAny(points)
}

func (m PointMap) Copy() PointMap {
	res := make(PointMap, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

func (m PointMap) PointsList(w, h int) (mtx [][]string) {
	res := make([][]string, h)
	for p, s := range m {
		if res[p.Y] == nil {
			res[p.Y] = make([]string, w)
		}
		res[p.Y][p.X] = s
	}
	return res
}

func (m PointMap) SetPoints(points [][]string) {
	for y, row := range points {
		for x, s := range row {
			if s != "" {
				m[Point{x, y}] = s
			}
		}
	}
}

func (m PointMap) Clear() {
	for p := range m {
		delete(m, p)
	}
}

func MergePoints(points []Point, color string, blocks PointMap) {
	for _, p := range points {
		blocks[p] = color
	}
}

func MergeTile(t *Tile, blocks PointMap) {
	MergePoints(t.Points(), Blocks[t.Typ()], blocks)
}
