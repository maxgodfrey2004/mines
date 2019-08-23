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
	"time"
)

const (
	// ChanceOfMine represents the chance (as a percentage) that any given grid cell will
	// contain a mine. Note that this value is approximate.
	ChanceOfMine = 20

	// MineRune represents a mine cell in char form.
	MineRune rune = 'M'
	// EmptyRune represents an empty cell in char form.
	EmptyRune rune = 'E'

	// EmptyRuneUser is the rune a user sees if they have checked a cell, but it contains nothing.
	EmptyRuneUser rune = '#'
	// UncheckedRuneUser is the rune a user sees if they have an unchecked cell.
	UncheckedRuneUser rune = '.'
)

// GridRow represents a row of the minesweeper grid.
type GridRow []rune

// GridType represents the minesweeper grid.
type GridType []GridRow

var (
	adjacentDeltas = []point{
		{Row: -1, Column: -1},
		{Row: -1, Column: 0},
		{Row: -1, Column: 1},
		{Row: 0, Column: -1},
		{Row: 0, Column: 1},
		{Row: 1, Column: -1},
		{Row: 1, Column: 0},
		{Row: 1, Column: 1},
	}
)

// makeGrid returns a GridType object containing cells being either empty or containing a mine.
func makeGrid(width, height int) GridType {
	grid := make(GridType, height)
	for i := 0; i < height; i++ {
		grid[i] = make(GridRow, width)
		for j := 0; j < width; j++ {
			grid[i][j] = EmptyRune
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	occupiedCells := 0
	var requiredCells int = (width * height) / (100 / ChanceOfMine)

	for occupiedCells < requiredCells {
		randRow := rand.Intn(height)
		randColumn := rand.Intn(width)
		if grid[randRow][randColumn] == EmptyRune {
			grid[randRow][randColumn] = MineRune
			occupiedCells++
		}
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if grid[i][j] != EmptyRune {
				continue
			}
			surroundingMines := 0
			for _, delta := range adjacentDeltas {
				newRow := i + delta.Row
				newColumn := j + delta.Column
				if newRow >= 0 && newRow < height && newColumn >= 0 && newColumn < width {
					if grid[newRow][newColumn] == MineRune {
						surroundingMines++
					}
				}
			}
			grid[i][j] = rune(int('0') + surroundingMines)
		}
	}

	return grid
}

// makeUserGrid returns a blank GridType object full of runes representing an unchecked cell.
func makeUserGrid(width, height int) GridType {
	grid := make(GridType, height)
	for i := 0; i < height; i++ {
		grid[i] = make(GridRow, width)
		for j := 0; j < width; j++ {
			grid[i][j] = UncheckedRuneUser
		}
	}

	return grid
}

// Game represents an instance of the Minesweeper game.
type game struct {
	grid          GridType      // The board on which the game is being played. The user does not see this board.
	userGrid      GridType      // The grid which the user sees when they play the game.
	selectedIndex point         // The current selected point in the grid.
	keypressChan  chan keypress // Incoming keyboard events for individual processing.
	Width         int           // The width of the game board.
	Height        int           // The height of the game board.
}

// New returns a new instance of the type game.
func New(width, height int) (g game) {
	g.grid = makeGrid(width, height)
	g.userGrid = makeUserGrid(width, height)
	g.Width = width
	g.Height = height
	return
}
