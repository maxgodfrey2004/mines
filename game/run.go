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
	"github.com/nsf/termbox-go"
)

// KeyEvent enumerates the different actions a user may take when using the application.
type KeyEvent int

const (
	// MoveUp represents an upwards movement of the cursor
	MoveUp KeyEvent = iota + 1
	// MoveDown represents a downwards movement of the cursor.
	MoveDown
	// MoveLeft represents a movement of the cursor to the left.
	MoveLeft
	// MoveRight represents a movement of the cursor to the right.
	MoveRight
	// Select represents the selection of a cell on the game board.
	Select
	// Flag represents the flagging of a grid cell
	Flag

	// Quit represents a quitting of the application.
	Quit
)

// keypress represents a keypress.
type keypress struct {
	EventType KeyEvent
	Key       termbox.Key
}

// point represents a coordinate on the grid.
type point struct {
	Row    int // The row at which the point is located.
	Column int // The column at which the point is located.
}

// listenForEvents listens for keypresses. Note that this is a blocking function so it must be
// called as a goroutine.
func (g *game) listenForEvents() {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				g.keypressChan <- keypress{EventType: MoveUp, Key: ev.Key}
			case termbox.KeyArrowDown:
				g.keypressChan <- keypress{EventType: MoveDown, Key: ev.Key}
			case termbox.KeyArrowLeft:
				g.keypressChan <- keypress{EventType: MoveLeft, Key: ev.Key}
			case termbox.KeyArrowRight:
				g.keypressChan <- keypress{EventType: MoveRight, Key: ev.Key}
			case termbox.KeyEnter:
				g.keypressChan <- keypress{EventType: Select, Key: ev.Key}
			default:
				switch ev.Ch {
				case rune('Q'), rune('q'):
					g.keypressChan <- keypress{EventType: Quit, Key: ev.Key}
				case rune('W'), rune('w'):
					g.keypressChan <- keypress{EventType: MoveUp, Key: ev.Key}
				case rune('S'), rune('s'):
					g.keypressChan <- keypress{EventType: MoveDown, Key: ev.Key}
				case rune('A'), rune('a'):
					g.keypressChan <- keypress{EventType: MoveLeft, Key: ev.Key}
				case rune('D'), rune('d'):
					g.keypressChan <- keypress{EventType: MoveRight, Key: ev.Key}
				case rune('F'), rune('f'):
					g.keypressChan <- keypress{EventType: Flag, Key: ev.Key}
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventInterrupt:
			return
		case termbox.EventResize:
			g.Render()
		}
	}
}

// moveCursor shifts the selected grid cell by a specified row and column delta. Note that the
// cursor is only ever shifted if the new specified location is within the grid. If it is not,
// the the cursor remains in the same place.
func (g *game) moveCursor(rowDelta, columnDelta int) {
	newRow := g.selectedIndex.Row + rowDelta
	newColumn := g.selectedIndex.Column + columnDelta

	if newRow >= 0 && newRow < g.Height && newColumn >= 0 && newColumn < g.Width {
		g.selectedIndex.Row = newRow
		g.selectedIndex.Column = newColumn
		g.Render()
	}
}

// flagCell places a flag on a given cell of the user's grid. Note that if the chosen cell is
// already flagged, then the flag is removed from that cell.
func (g *game) flagCell(row, column int) {
	if g.GameOver {
		return
	}
	if g.userGrid[row][column] == UncheckedRuneUser && g.flaggedCells < g.maxFlags {
		g.userGrid[row][column] = FlaggedRuneUser
		g.flaggedCells++
	} else if g.userGrid[row][column] == FlaggedRuneUser {
		g.userGrid[row][column] = UncheckedRuneUser
		g.flaggedCells--
	}
	g.Render()
	g.finishGameIfWon()
}

// selectAllMines reveals the positions of all mines on the game board.
func (g *game) selectAllMines() {
	for _, minePos := range g.mines {
		g.showCell(minePos.Row, minePos.Column)
	}
}

// selectCell selects a given row and column of the grid. It then reveals all surrounding cells
// according to the rules of the game. If the current cell is a mine, it is game over!
func (g *game) selectCell(row, column int) {
	if g.firstMove {
		g.handleFirstMove(row, column)
	}
	if g.GameOver || g.userGrid[row][column] == FlaggedRuneUser {
		return
	}
	if g.grid[row][column] == rune('0') {
		g.showCell(row, column)
		g.selectFlood(row, column)
	}
	if g.grid[row][column] == MineRune {
		g.GameOver = true
		g.selectAllMines()
	}
	if int(g.grid[row][column]) >= int('1') && int(g.grid[row][column]) <= int('9') {
		g.showCell(row, column)
	}
	g.Render()
}

// selectFlood selects every cell surrounding an empty cell. This function is designed to be called
// when a grid cell with no surrounding mines is selected.
func (g *game) selectFlood(row, column int) {
	selectionQueue := make([]point, 1)
	selectionQueue[0] = point{Row: row, Column: column}

	for len(selectionQueue) > 0 {
		currentPoint := selectionQueue[0]
		selectionQueue = selectionQueue[1:]

		for _, deltaPoint := range adjacentDeltas {
			newPoint := point{
				Row:    currentPoint.Row + deltaPoint.Row,
				Column: currentPoint.Column + deltaPoint.Column,
			}
			if g.inGrid(newPoint.Row, newPoint.Column) {
				if g.userGrid[newPoint.Row][newPoint.Column] == UncheckedRuneUser {
					if g.grid[newPoint.Row][newPoint.Column] != MineRune {
						g.showCell(newPoint.Row, newPoint.Column)
						if g.userGrid[newPoint.Row][newPoint.Column] == EmptyRuneUser {
							selectionQueue = append(selectionQueue, newPoint)
						}
					}
				}
			}
		}
	}
}

// showCell displays a previously unchecked cell on the grid which the user sees.
func (g *game) showCell(row, column int) {
	switch g.grid[row][column] {
	case rune('0'):
		g.userGrid[row][column] = EmptyRuneUser
	default:
		g.userGrid[row][column] = g.grid[row][column]
	}
	g.shownCells++

	g.finishGameIfWon()
}

// finishGameIfWon makes the grid fully green to indicate success if the user has filled it in and
// flagged every mine correctly.
func (g *game) finishGameIfWon() {
	if g.shownCells+g.flaggedCells == g.Width*g.Height && !g.GameOver {
		g.GameOver = true
		for i := 0; i < g.Height; i++ {
			for j := 0; j < g.Width; j++ {
				g.userGrid[i][j] = GameWonRuneUser
			}
		}
		g.Render()
	}
}

// Run starts the game.
func (g *game) Run() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	g.keypressChan = make(chan keypress)
	g.selectedIndex = point{Row: 0, Column: 0}

	g.Render()
	go g.listenForEvents()

	for {
		select {
		case ev := <-g.keypressChan:
			switch ev.EventType {
			case MoveUp:
				g.moveCursor(-1, 0)
			case MoveDown:
				g.moveCursor(1, 0)
			case MoveLeft:
				g.moveCursor(0, -1)
			case MoveRight:
				g.moveCursor(0, 1)
			case Quit:
				termbox.Close()
				return
			case Select:
				g.selectCell(g.selectedIndex.Row, g.selectedIndex.Column)
			case Flag:
				g.flagCell(g.selectedIndex.Row, g.selectedIndex.Column)
			}
		}
	}
}
