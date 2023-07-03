package main

import (
	"testing"
)

func TestNewTileCoordinateBound(t *testing.T) {
	bound := NewTileCoordinateBound(
		WithTileSize(256),
		WithLevel(15),
		WithBounds(27380, 13434, 27380, 13434))
	bound.TileBounds()
}
