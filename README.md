# mines [![Go Report Card](https://goreportcard.com/badge/github.com/maxgodfrey2004/mines)](https://goreportcard.com/report/github.com/maxgodfrey2004/mines)  [![Maintainability](https://api.codeclimate.com/v1/badges/259fa9fb2c430a3e0300/maintainability)](https://codeclimate.com/github/maxgodfrey2004/mines/maintainability)

A terminal-based implementation of Minesweeper in Go.

![mines_demo_better](https://user-images.githubusercontent.com/34620214/63650957-eb5e1c00-c782-11e9-94f3-6f5cff5550c0.png)

*An example of a medium level game*

## Contents

1. [Building and Installing](#Building-and-Installing)
    - [Installing](#Installing)
    - [Building](#Building)
2. [Controls](#Controls)
3. [Game Reference](#Game-Reference)
    - [How to play](#How-to-play)
    - [Game modes](#Game-modes)
4. [License](#License)

This project is licensed under the Apache-2](#License

This project is licensed under the Apache-2)

## Building and Installing

### Installing

mines has a fairly minimal set of dependencies:
  - [Go](https://golang.org/doc/install) 1.10
  - termbox-go (github.com/nsf/termbox-go)

To install the application, run the following command:

```bash
go get -u github.com/maxgodfrey2004/mines
```

### Building

In future, a Makefile will be added here so that building the project is easier. However, in the meantime, you will have to run the build commands manually.

1. Change your working directory to the repository root (this will be located somewhere in your GOPATH). If you are not sure where your GOPATH is, just execute:

```bash
go env GOPATH
```

2. Build the project by running the following:

```bash
go build -o mines main.go
```

### Running

Now you can run your new binary. If you require any assistance pertinent to the usage of the binary, invoke it with the `--help` flag for more information.

### A brief note concerning command line arguments

Passing flags to the binary is syntactically easier than you may expect. The Go [flag](https://golang.org/pkg/flag/) package allows for many different ways to pass command line argments.

All of the following are valid:

```bash
./mines -difficulty medium
./mines -difficulty=medium
./mines --difficulty medium
./mines --difficulty=medium
```

You can also put quotes around the argument's value if you are that way inclined.

## Controls

| Key                     | Functionality                                          |
|-------------------------|--------------------------------------------------------|
| `Arrow Up`, `W`, `w`    | Moves the selected cell up by one.                     |
| `Arrow Down`, `S`, `s`  | Moves the selected cell down by one.                   |
| `Arrow Left`, `A`, `a`  | Moves the selected cell left by one.                   |
| `Arrow Right`, `D`, `d` | Moves the selected cell right by one.                  |
| `F`, `f`                | Flags the selected cell as a mine.                     |
| `Q`, `q`                | Quits the application.                                 |
| `Return`                | Selects the current cell and displays its contents.    |

## Game Reference

### How to play

Minesweeper is played on a grid. There are two different types of grid cell, these being numbered cells and mines. If you select a mine, you lose and the game is over.

The number in a numbered cell represents the amount of mines directly adjacent to that cell. For an example of adjacent cells, have a look at the following grid - cells marked `A` are considered adjacent to the cell marked `X`.

```
..........
..........
.....AAA..
.....AXA..
.....AAA..
..........
```

You start off with an empty grid, where the value of every cell is hidden. The goal is to find the location of every mine on the grid (thus "sweeping" the grid of mines, hence the name of the game) without actually selecting a mine. Given the information presented to you by various numbered cells, it is possible to make informed decisions as to where mines are.

If there is a cell you are adamant contains a mine, do not select it (if you do, then you lose the game and all your progress!). Instead, flag it. Note that you are not able to select a cell which you have already flagged. This stops you from accidentally selecting that cell and losing the game.

### Game modes

mines can be played in three modes: "Easy", "Medium" and "Hard".

The difficulty of each game can be set through command line flags when the application is invoked. If the difficulty of a game is not set, then the application defaults to medium difficulty.

The specifications for each game mode are as follows:

- Easy

  |            |             |
  |------------|-------------|
  | Width      | 8 cells     |
  | Height     | 8 cells     |
  | Mines      | 10          |

- Medium

  |            |             |
  |------------|-------------|
  | Width      | 16 cells    |
  | Height     | 16 cells    |
  | Mines      | 40          |

- Hard

  |            |             |
  |------------|-------------|
  | Width      | 30 cells    |
  | Height     | 16 cells    |
  | Mines      | 99          |

## License

This project is licensed under the Apache License 2.0. For more information, read the LICENCE file in the project's root directory. Alternatively, you can [read it on Github](https://github.com/maxgodfrey2004/mines/blob/master/LICENSE).
