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

// renderString draws a specified string on the terminal, starting at a given x and y coordinate.
// The parameter `str` represents the string to be drawn.
func renderString(x, y int, str []rune) {
	for i := 0; i < len(str); i++ {
		termbox.SetCell(x+i, y, rune(str[i]), termbox.ColorDefault, termbox.ColorDefault)
	}
}

// Render draws the game's grid on the terminal screen.
func (g *game) Render() {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	// Make sure that the terminal's back buffer is large enough for us to draw
	// the full grid.
	tWidth, tHeight := termbox.Size()

	if g.Width > tWidth || g.Height > tHeight {
		renderString(0, 0, []rune("Terminal is too small. Please resize."))
		return
	}

	for i, rowString := range g.grid {
		renderString(0, i, rowString)

		termbox.SetCursor(g.selectedIndex.Column, g.selectedIndex.Row)

		// for j := 0; j < len(rowString); j++ {
		// 	fgColor := termbox.ColorDefault
		// 	bgColor := termbox.ColorDefault
		// }
	}

	termbox.Flush()
}
