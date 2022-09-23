package gotris

type (
	Typ  string
	Dir  int
	Spin int
)

const (
	// TODO: use the names defined by the community
	TypB Typ = "B"
	TypI Typ = "I"
	TypL Typ = "L"
	TypJ Typ = "J"
	TypT Typ = "T"
	TypS Typ = "S"
	TypZ Typ = "Z"
)

const (
	DirUp    Dir = 0
	DirRight Dir = 1
	DirDown  Dir = 2
	DirLeft  Dir = 3
)

const (
	SpinLeft  Spin = 0
	SpinRight Spin = 1
)

type Point struct {
	X int
	Y int
}

type Tile struct {
	Type   Typ
	Points []Point
}

func NewTile(typ Typ, x, y int) *Tile {
	return &Tile{
		Type:   typ,
		Points: PointsForType(typ, x, y),
	}
}

func RandomTileType() Typ {
	return TypB
}

func PointsForType(typ Typ, x, y int) []Point {
	return []Point{{x, y}}
}
