package game

import (
	"math/rand"

	"github.com/ubntc/go/games/gotris/common/geometry"
)

const (
	TypB geometry.Typ = "B"
	TypI geometry.Typ = "I"
	TypL geometry.Typ = "L"
	TypJ geometry.Typ = "J"
	TypT geometry.Typ = "T"
	TypS geometry.Typ = "S"
	TypZ geometry.Typ = "Z"
)

var Tiles = []struct {
	Typ geometry.Typ
	// Drawings defines how the tile is drawn using a small DSL:
	//  'x' = draw + move right
	//  'c' = draw + move right + use as center point
	//  ' ' = move right
	//  ',' = next row
	Drawings []string
}{
	// All orientations and center blocks in 7 lines of code! 🤯
	{TypB, []string{"xc,xx", "xc,xx", "xc,xx", "xc,xx"}},
	{TypI, []string{"x,c,x,x", "xcxx", "x,x,c,x", "xxcx"}},
	{TypL, []string{"x,c,xx", "xcx,x", "xx, c, x", "  x,xcx"}},
	{TypJ, []string{" x, c,xx", "x,xcx", "xx,c,x", "xcx,  x"}},
	{TypT, []string{"xcx, x", " x,xc, x", " x,xcx", "x,cx,x"}},
	{TypS, []string{" cx,xx", "x,xc, x", " xx,xc", "x,cx, x"}},
	{TypZ, []string{"xc, xx", " x,xc,x", "xx, cx", " x,cx,x"}},
}

var drawings geometry.Drawings

func init() {
	drawings = make(geometry.Drawings, len(Tiles))
	for _, spec := range Tiles {
		drawings[spec.Typ] = spec.Drawings
	}
}

func RandomTile(rng *rand.Rand) *geometry.Tile {
	i := rng.Int() % len(Tiles)
	typ := Tiles[i].Typ
	return geometry.NewTile(typ, 0, 0, drawings)
}
