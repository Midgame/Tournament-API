package main

import (
	"testing"
)

func TestInit(t *testing.T) {
	if len(knownBots) != 0 {
		t.Errorf("Known bots wasn't initialized to 0 properly")
	}
	if grid.Height != GRID_HEIGHT {
		t.Errorf("Height not initialized properly for grid")
	}
	if grid.Width != GRID_WIDTH {
		t.Errorf("Width not initialized properly for grid")
	}
	if len(grid.Entities) != 0 {
		t.Errorf("Grid entities not initialized properly")
	}
}
