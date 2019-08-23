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

// makeGrid propogates the game grid with mines, and precomputes the amount of surrounding mines
// for every cell that does not contain a mine.
func (g *game) makeGrid() {
	g.grid = make(GridType, g.Height)
	for i := 0; i < g.Height; i++ {
		g.grid[i] = make(GridRow, g.Width)
		for j := 0; j < g.Width; j++ {
			g.grid[i][j] = EmptyRune
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	occupiedCells := 0
	var requiredCells int = (g.Width * g.Height) / (100 / ChanceOfMine)

	for occupiedCells < requiredCells {
		randRow := rand.Intn(g.Height)
		randColumn := rand.Intn(g.Width)
		if g.grid[randRow][randColumn] == EmptyRune {
			g.grid[randRow][randColumn] = MineRune
			g.mines = append(g.mines, point{Row: randRow, Column: randColumn})
			occupiedCells++
		}
	}

	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			if g.grid[i][j] != EmptyRune {
				continue
			}
			surroundingMines := 0
			for _, delta := range adjacentDeltas {
				newRow := i + delta.Row
				newColumn := j + delta.Column
				if newRow >= 0 && newRow < g.Height && newColumn >= 0 && newColumn < g.Width {
					if g.grid[newRow][newColumn] == MineRune {
						surroundingMines++
					}
				}
			}
			g.grid[i][j] = rune(int('0') + surroundingMines)
		}
	}
}

// makeUserGrid propogates the grid which the user sees with the rune representing an unseen cell.
func (g *game) makeUserGrid() {
	g.userGrid = make(GridType, g.Height)
	for i := 0; i < g.Height; i++ {
		g.userGrid[i] = make(GridRow, g.Width)
		for j := 0; j < g.Width; j++ {
			g.userGrid[i][j] = UncheckedRuneUser
		}
	}
}

// Game represents an instance of the Minesweeper game.
type game struct {
	grid          GridType      // The board on which the game is being played. The user does not see this board.
	userGrid      GridType      // The grid which the user sees when they play the game.
	selectedIndex point         // The current selected point in the grid.
	mines         []point       // The locations of mines on the grid.
	keypressChan  chan keypress // Incoming keyboard events for individual processing.
	GameOver      bool          // Represents whether or not the user has lost the game.
	Width         int           // The width of the game board.
	Height        int           // The height of the game board.
}

// New returns a new instance of the type game.
func New(width, height int) (g game) {
	g.Width = width
	g.Height = height
	g.makeGrid()
	g.makeUserGrid()
	return
}
