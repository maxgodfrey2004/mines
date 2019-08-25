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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maxgodfrey2004/mines/game"
)

const (
	DifficultyDefault = "medium"
	DifficultyDesc    = "The difficulty of the game [\"easy\"|\"medium\"|\"hard\"]"
)

var (
	gameDifficulty string
	gameWidth      int
	gameHeight     int
	gameNumMines   int
)

// flagInit initialises all command line flags
func flagInit() {
	flag.StringVar(&gameDifficulty, "difficulty", DifficultyDefault, DifficultyDesc)
	flag.Parse()
}

func checkFlags() {
	switch gameDifficulty {
	case "easy":
		gameWidth = 8
		gameHeight = 8
		gameNumMines = 10
	case "medium":
		gameWidth = 16
		gameHeight = 16
		gameNumMines = 40
	case "hard":
		gameWidth = 30
		gameHeight = 16
		gameNumMines = 99
	default:
		fmt.Fprintln(os.Stderr, "game difficulty must be [\"easy\"|\"medium\"|\"hard\"]")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	flagInit()
	checkFlags()

	minesweeper := game.New(gameWidth, gameHeight, gameNumMines)
	minesweeper.Run()
}
