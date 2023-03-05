package geometry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	geom "github.com/ubntc/go/games/gotris/common/geometry"
)

func TestPoint(t *testing.T) {
	assert := assert.New(t)

	a := geom.Point{X: 1, Y: 1}
	b := geom.Point{X: 1, Y: 1}

	assert.True(a == b, "ensure points can be compared")
}
