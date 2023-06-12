package game

import (
	"fmt"
	"time"
)

func (g *Game) Dump() {
	fmt.Println("Blocks", g.Board)
	fmt.Println("CurrentTile", g.CurrentTile)
	fmt.Println("NextTile", g.NextTile)
}

// nolint
func hint[K any](v ...K) {
	fmt.Printf("%v\n", v)
	time.Sleep(time.Second)
}
