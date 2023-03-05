package geometry

// PointMap stores strings in a 2D map with Points(x,y) as the key.
type PointMap map[Point]string

// ContainsAny returns whether or not the map has a value defined
// for at least one of the given points.
func (m PointMap) ContainsAny(points []Point) bool {
	for _, p := range points {
		if _, ok := m[p]; ok {
			return true
		}
	}
	return false
}

// Collides checks whether or not the given Points, if moved in the given
// direction, would have a value defined in the map.
func (m PointMap) Collides(points []Point, dir Dir) bool {
	points = OffsetPointsDir(points, dir)
	return m.ContainsAny(points)
}

// Copy returns a deep copy of the map.
func (m PointMap) Copy() PointMap {
	res := make(PointMap, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

// PointsList returns the map as a 2D list (slice of string slices).
// The given width w and height h define the dimensions of the slice.
// The slice is then filled with any values stored in the map at [0:w],[0:h].
// Any values outside of this range are not added to the slice.
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

// SetPoints sets the given string values in the map, using the corresponding
// slice indices as coordinates.
func (m PointMap) SetPoints(points [][]string) {
	for y, row := range points {
		for x, s := range row {
			if s != "" {
				m[Point{x, y}] = s
			}
		}
	}
}

// Clear deletes all points in the map.
func (m PointMap) Clear() {
	for p := range m {
		delete(m, p)
	}
}

// Set sets the given string for the given point in the map and
// returns whether or not the value was already present at the given point.
func (m PointMap) Set(point Point, value string) bool {
	if m[point] == value {
		return true
	}
	m[point] = value
	return false
}

// SetAll sets the given string for all given points in the map and
// retuns the points that provided a new value to the map.
func (m PointMap) SetAll(points []Point, value string) (res []Point) {
	for _, p := range points {
		if m.Set(p, value) {
			res = append(res, p)
		}
	}
	return
}
