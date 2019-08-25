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
	// MineRune represents a mine cell in char form.
	MineRune rune = 'M'
	// EmptyRune represents an empty cell in char form.
	EmptyRune rune = 'E'

	// EmptyRuneUser is the rune a user sees if they have checked a cell, but it contains nothing.
	EmptyRuneUser rune = '#'
	// GameWonRuneUser is the rune displayed on every grid cell when/if the user wins the game.
	GameWonRuneUser rune = '*'
	// FlaggedRuneUser is the rune a user sees when they have flagged a cell they believe contains a mine.
	FlaggedRuneUser rune = 'F'
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

// game represents an instance of the Minesweeper game.
type game struct {
	flaggedCells  int           // The number of cells which the user has flagged.
	shownCells    int           // The number of cells which have been shown to the user.
	maxFlags      int           // The maximum number of cells which the user can place flags on.
	grid          GridType      // The board on which the game is being played. The user does not see this board.
	userGrid      GridType      // The grid which the user sees when they play the game.
	selectedIndex point         // The current selected point in the grid.
	numMines      int           // The amount of mines to be placed on the grid.
	mines         []point       // The locations of mines on the grid.
	keypressChan  chan keypress // Incoming keyboard events for individual processing.
	GameOver      bool          // Represents whether or not the user has lost the game.
	Width         int           // The width of the game board.
	Height        int           // The height of the game board.
}

// inGrid returns whether or not a given location is within the game grid.
func (g *game) inGrid(row, column int) bool {
	return row >= 0 && row < g.Height && column >= 0 && column < g.Width
}

// makeGrid propagates the game grid with mines, and precomputes the amount of surrounding mines
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
	g.maxFlags = g.numMines

	for occupiedCells < g.maxFlags {
		randRow := rand.Intn(g.Height)
		randColumn := rand.Intn(g.Width)
		if g.grid[randRow][randColumn] == EmptyRune {
			g.grid[randRow][randColumn] = MineRune
			g.mines = append(g.mines, point{Row: randRow, Column: randColumn})
			occupiedCells++
		}
	}

	g.precomputeSurroundingMines()
}

// makeUserGrid propagates the grid which the user sees with the rune representing an unseen cell.
func (g *game) makeUserGrid() {
	g.userGrid = make(GridType, g.Height)
	for i := 0; i < g.Height; i++ {
		g.userGrid[i] = make(GridRow, g.Width)
		for j := 0; j < g.Width; j++ {
			g.userGrid[i][j] = UncheckedRuneUser
		}
	}
}

// precomputeSurroundingMines precomputes the amount of surrounding mines for every cell that does
// not contain a mine. This can be expressed as a single rune because the maximum amount of mines
// surrounding a cell is nine, which is expressed as `rune('9')`.
func (g *game) precomputeSurroundingMines() {
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			if g.grid[i][j] != EmptyRune {
				continue
			}
			surroundingMines := 0
			for _, delta := range adjacentDeltas {
				newRow := i + delta.Row
				newColumn := j + delta.Column
				if g.inGrid(newRow, newColumn) {
					if g.grid[newRow][newColumn] == MineRune {
						surroundingMines++
					}
				}
			}
			g.grid[i][j] = rune(int('0') + surroundingMines)
		}
	}
}

// New returns a new instance of the type game.
func New(width, height, numMines int) (g game) {
	g.flaggedCells = 0
	g.shownCells = 0
	g.numMines = numMines
	g.Width = width
	g.Height = height
	g.makeGrid()
	g.makeUserGrid()
	return
}
