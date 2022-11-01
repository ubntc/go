# Gotris
```
  　🟪🟪　╔═╗╔═╗╔╦╗╦═╗╦╔═╗　🟦🟦　　
  🟪🟪　　║ ╦║ ║ ║ ╠╦╝║╚═╗　　🟦🟦　
  　　　　╚═╝╚═╝ ╩ ╩╚═╩╚═╝　　　　　
```
Gotris is an experimental implementation of a the famous block dropping game.

The goal of this project is to learn how to implement all aspects of such a game
using only the Go standard library without much need for external dependencies.
In some places dependecies may be pulled in to avoid too much custom code.

```
┌一一一一一一一一┐
│　　　　　　　　│  ＮＥＸＴ
│　　　　　　　　│
│　　　　　　　　│  🟫🟫　　
│　　　　🟥　　　│  🟫🟫　　
│　　　　🟥　　　│
│　　　　🟥🟥　　│
│　　　　　　　　│  ＳＣＯＲＥ
│　　　　　　　　│
│　　　　　　　　│  0
│　　　　　　　　│
│🟨　　　　　　　│  ＬＥＶＥＬ
│🟨🟧　　🟦　　　│
│🟨🟧🟧🟧🟦🟦🟨　│  1000
│🟨🟪🟩🟩🟩🟦🟨　│
│🟪🟪🟩🟩🟫🟫🟨　│
│🟪🟩🟩🟩🟫🟫🟨　│
│￣￣￣￣￣￣￣￣│
└一一一一一一一一┘
```

The game is very playable and runs on Linux and MacOS.

## Running the game

Git clone this repo and run it with `go run cmd/gotris/gotris.go` or directly install it using:
```
go install github.com/ubntc/go/games/gotris/cmd/gotris@latest
```
And then just run `gotris`.

## Features

 1. 🌈 Colors! (requires unicode terminal)
 2. 🚀 Low lag and fast input! (for me it is really fun to play)
 3. 🪄 Smart rotation! (allow rotation on edge of the board)
 4. 🔢 Many keys mapped! (press H or ? to see controls)
 5. 🫣 Preview of next tile! (without would not be fun)
 6. 👨‍🎨 ASCII/Unicode Art (mindblowing title and help screen)

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

## [Screenshots](Screenshots.md)

## [Bugs](Bugs.md)

## Tested On

* Linux Rasberry Pi 4 (32bit, go 1.19, needs fonts-noto-color-emoji package)
* MacOS iTerm/Appleterm and the VSCode terminal (M1 CPU, go 1.18)

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
1. Tests and playing with various options to "clear" the terminal
1. Random seed and make non-fixed see the default.
1. Playing with testability of the terminal (capture input, get width in tests)
1. While doing this, make Terminal a class with stdout configurable
1. Function docs for many public methods.
1. More key mappings (forgot to setup WASD 😅)
1. Title screen and game over screen
1. Experiments with more full-width anmd half-width box drawing
1. Use nice full-width letters for true blocky look
1. More text art and a help screen
1. More files to subpackages to unclutter the game package
1. Make make platform and input structs of the game to avoid hacky closures
1. More explicit key handling (do not just send runes)
1. Create shared input package to resolve dependcy cycle
1. Grab key modifiers and use them for something cool (shift game screen)
1. More comments and removal of dead code
1. Some core method renaming and signature changes gave some ideas for new features
1. New platform feature to print (debug) messages in a nice place
1. Fixed broken tests that did not respect the new way of input handling
1. Changed frameArts to be more usable and exchangeable
1. Options menu with render options
1. In-game rendermode switch using key 1 to 9
1. Move static modes code to a mode manager
1. Awesome menu system and refactored screen drawing system
1. Integrate a real GUI platform: fyne.io v2
1. Realized that strigs with escape codes look terrible on as fyne Canvas text
1. Implemented changes to support platform switch via CLI
1. Dummy platform to as copy&paste for future experiments and tests
1. Moved capture on/off to game config
1. Refactored menu system
  * Made it generic on the Game side and concrete on the Platform side.
  * Mapped all menus/scenes to the concrete textscenes in the text platform.
  * moves scene loops to separate files names scene-num-name.go
  * This change went unexpectedly smooth and was done quickly. Go is really awesome for heavy code reshaping! Also putting things in subpackages and naming all magic strings from the beginning helped a  lot.
1. Define fyne gray scale colors.
1. Not happy with how Options are set in the platform. Refactored it again. Seems to be getting in shape now. Refactoring is still fun!
1. Moved all scenes to separate files to see what event loop logic is actually needed to manage a scene. This led to the extraction of some reusabe handlers, such as the global key handler and the options handler. It also allowed moving out the options index logic to a central place.