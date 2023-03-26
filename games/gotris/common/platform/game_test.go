package platform_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ubntc/go/games/gotris/common/geometry"
	"github.com/ubntc/go/games/gotris/common/platform"
)

func TestMarshallingGame(t *testing.T) {
	g1 := platform.Game{
		Rules: platform.Rules{
			TickTime:    1,
			SpeedStep:   1,
			MaxSteps:    1,
			BoardSize:   *geometry.NewDim(10, 10),
			PreviewSize: *geometry.NewDim(1, 1),
			Seed:        1,
		},
		Score:    1,
		Speed:    time.Millisecond,
		BoardPos: *geometry.NewDim(1, 1),
		NextTile: geometry.NewTile("L", 1, 1, geometry.Drawings{
			"L": []string{"x,c,xx"},
		}),
		CurrentTile: geometry.NewTile("I", 1, 1, geometry.Drawings{
			"I": []string{"x,x,c,x"},
		}),
		Board: geometry.PointMap{
			*geometry.NewPoint(1, 1): "x",
		},
	}

	data, err := json.Marshal(g1)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	g2 := platform.Game{}

	err = json.Unmarshal(data, &g2)
	assert.NoError(t, err)

	assert.Equal(t, g1, g2)
}
