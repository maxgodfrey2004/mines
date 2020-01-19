// Copyright 2020 Max Godfrey
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

// Various color constants to make rendering easier.
const (
	bgColorDefault = termbox.ColorDefault
	fgColorDefault = termbox.ColorBlue

	fgColorFlag = termbox.ColorYellow
	bgColorFlag = bgColorDefault

	fgColorMine = termbox.ColorRed
	bgColorMine = bgColorDefault

	fgColorGameWon = termbox.ColorGreen
	bgColorGameWon = bgColorDefault
)

// renderString draws a specified string on the terminal, starting at a given x and y coordinate.
// The parameter `str` represents the string to be drawn.
func renderString(x, y int, str []rune) {
	for i := 0; i < len(str); i++ {
		termbox.SetCell(x+i, y, str[i], fgColorDefault, bgColorDefault)
	}
}

// resolveColors returns the preset foreground and background colors for a given rune.
func resolveColors(gridCell rune) (fgColor, bgColor termbox.Attribute) {
	switch gridCell {
	case MineRune:
		return fgColorMine, bgColorMine
	case FlaggedRuneUser:
		return fgColorFlag, bgColorFlag
	case GameWonRuneUser:
		return fgColorGameWon, bgColorGameWon
	}
	return fgColorDefault, bgColorDefault
}

// Render draws the game's grid on the terminal screen.
func (g *game) Render() {
	if err := termbox.Clear(fgColorDefault, bgColorDefault); err != nil {
		panic(err)
	}

	// Make sure that the terminal's back buffer is large enough for us to draw
	// the full grid.
	tWidth, tHeight := termbox.Size()
	if g.Width > tWidth || g.Height > tHeight {
		renderString(0, 0, []rune("Terminal is too small. Please resize."))
		return
	}

	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			fgColor, bgColor := resolveColors(g.userGrid[i][j])
			if i == g.selectedIndex.Row && j == g.selectedIndex.Column {
				// Invert the cell's colors to indicate to the user that the current
				// cell has been selected.
				fgColor, bgColor = bgColor, fgColor
			}
			termbox.SetCell(j, i, g.userGrid[i][j], fgColor, bgColor)
		}
	}

	termbox.Flush()
}
