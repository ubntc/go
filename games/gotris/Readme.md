# Gotris

Gotris is an experimental implementation of a the famous block dropping game.

The goal of this project is to learn how to implement all aspects of such a game
using only the Go standard library without much need for external dependencies.
In some places dependecies may be pulled in to avoid too much custom code.

```
⎛﹋﹋﹋﹋﹋﹋﹋﹋⎞
│　　　　　　　　│  NEXT
│　　　　　　　　│  　　　　　
│　　　　　　　　│  🟪🟪　　　
│　　　　　　　　│  　🟪🟪　　
│　　🟩🟩🟩　　　│  　　　　　
│　　　🟩　　　　│
│　　　　　　　　│  SCORE
│　　　　　　　　│
│　　　　　　　　│  0
│　　　　　　　　│
│　　　　　　🟦　│  Speed
│　　　　　🟪🟦🟦│
│　🟧🟧　🟪🟪🟦🟦│  1000 ms
│　🟧🟫🟫🟪🟩🟦🟦│
│　🟧🟫🟫🟩🟩🟩🟦│
│﹋﹋﹋﹋﹋﹋﹋﹋│
⎝﹏﹏﹏﹏﹏﹏﹏﹏⎠
```

The game is very playable and runs on Linux and MacOS.

Run it with:

```
go run cmd/main.go
```

## Features

 1. 🌈 Colors! (requires unicode terminal)
 2. 🚀 Low lag and fast input! (for me it is really fun to play)
 3. 🪄 Smart rotation! (allow rotation on edge of the board)
 4. 🔢 All keys mapped! (so the game is fun to play in any country 🤞)
 5. 🫣 Preview of next tile! (without would not be fun)

## Missing Features

* Sound
* Official Scoring Rules
* Scoreboard
* Alternative Input
* Multiplayer on one screen
* Multiplayer across the network
* Top-Score Animations
* Non-Standard Tiles
* In-Game Menus and Options

## Tested On

* Linux Rasberry Pi 4 (32bit, go 1.19, needs fonts-noto-color-emoji package)
* MacOS in iTerm and the VSCode terminal (M1 CPU, go 1.18)

## Implementation History

This was intended more as a fun coding trip and not a full fledged project.
Therefore, here is my journey of things implemented.

  1. Basic game entities and game loop
  1. Tile movement and collision
  1. Basic terminal rendering
  1. User input (via stdin)
  1. Tile rotation (which triggerd quite some refactoring)
  1. Smart tile drawing using a short string to define how to draw a tile (instead of slices of numbers)
  1. Drawing specs for orientations with a well-defined center (it is more complicated than you may think)
  1. Smart rotation (move tile left or right to resolve blocked rotation)
  1. Scoring (only then the game was playable and fun 🥳🎊🎉)
  1. Next tile preview
  1. Rendering components (for preview and score box)
  1. Frames around view boxes
  1. Better speed adjustment
  1. Separation of concers: dissolve code of big main.go into the packages
  1. Make game package independent of the renderer and terminal (added "platform" notion)
  1. Adding tests and playing with various options to "clear" the terminal
  1. Add random seed and make non-fixed see the default.
  1. Playing with testability of the terminal (capture input, get width in tests)
  1. While doing this, make Terminal a class with stdout configurable
  1. Add function docs for many public methods.
  1. Add more key maopings (forgot to setup WASD 😅)
  1. Add title screen and game over screen
