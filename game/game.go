// Copyright 2019 Max Godfrey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package game

import (
	"math/rand"
)

const (
	// ChanceOfMine represents the chance (in decimal form) that any given grid cell will
	// contain a mine.
	ChanceOfMine = 0.2

	// MineRune represents a mine cell in char form.
	MineRune = 'M'
	// EmptyRune represents an empty cell in char form.
	EmptyRune = 'E'
)

type GridRow []rune
type GridType []GridRow

// makeGrid returns a GridType object containing cells being either empty or containing a mine.
func makeGrid(width, height int) GridType {
	grid := make(GridType, height)
	for i := 0; i < height; i++ {
		grid[i] = make(GridRow, width)
		for j := 0; j < width; j++ {
			if rand.Float64() < ChanceOfMine {
				grid[i][j] = MineRune
			} else {
				grid[i][j] = EmptyRune
			}
		}
	}
	return grid
}

// Game represents an instance of the Minesweeper game.
type game struct {
	grid          GridType      // The board on which the game is being played.
	selectedIndex point         // The current selected point in the grid.
	keypressChan  chan keypress // Incoming keyboard events for individual processing.
	Width         int           // The width of the game board.
	Height        int           // The height of the game board.
}

// New returns a new instance of the type game.
func New(width, height int) (g game) {
	g.grid = makeGrid(width, height)
	g.Width = width
	g.Height = height
	return
}
