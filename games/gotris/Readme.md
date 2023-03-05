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

Git clone this repo and run it with `go run gotris.go` or directly install it using:
```
go install github.com/ubntc/go/games/gotris@latest
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
* Linux x86 (64bit, go 1.20, Manjaro)
* MacOS iTerm/Appleterm and the VSCode terminal (M1 CPU, go 1.18)

## Implementation History

This was intended more as a fun coding trip and not a full fledged project.
Therefore, here is my journey of things implemented.

* Basic game entities and game loop
* Tile movement and collision
* Basic terminal rendering
* User input (via stdin)
* Tile rotation (which triggerd quite some refactoring)
* Smart tile drawing using a short string to define how to draw a tile (instead of slices of numbers)
* Drawing specs for orientations with a well-defined center (it is more complicated than you may think)
* Smart rotation (move tile left or right to resolve blocked rotation)
* Scoring (only then the game was playable and fun 🥳🎊🎉)
* Next tile preview
* Rendering components (for preview and score box)
* Frames around view boxes
* Better speed adjustment
* Separation of concers: dissolve code of big main.go into the packages
* Make game package independent of the renderer and terminal (added "platform" notion)
* Tests and playing with various options to "clear" the terminal
* Random seed and make non-fixed see the default.
* Playing with testability of the terminal (capture input, get width in tests)
* While doing this, make Terminal a class with stdout configurable
* Function docs for many public methods.
* More key mappings (forgot to setup WASD 😅)
* Title screen and game over screen
* Experiments with more full-width anmd half-width box drawing
* Use nice full-width letters for true blocky look
* More text art and a help screen
* More files to subpackages to unclutter the game package
* Make make platform and input structs of the game to avoid hacky closures
* More explicit key handling (do not just send runes)
* Create shared input package to resolve dependcy cycle
* Grab key modifiers and use them for something cool (shift game screen)
* More comments and removal of dead code
* Some core method renaming and signature changes gave some ideas for new features
* New platform feature to print (debug) messages in a nice place
* Fixed broken tests that did not respect the new way of input handling
* Changed frameArts to be more usable and exchangeable
* Options menu with render options
* In-game rendermode switch using key 1 to 9
* Move static modes code to a mode manager
* Awesome menu system and refactored screen drawing system
* Integrate a real GUI platform: fyne.io v2
* Realized that strigs with escape codes look terrible on as fyne Canvas text
* Implemented changes to support platform switch via CLI
* Dummy platform to as copy&paste for future experiments and tests
* Moved capture on/off to game config
* Refactored menu system
  * Made it generic on the Game side and concrete on the Platform side.
  * Mapped all menus/scenes to the concrete textscenes in the text platform.
  * moves scene loops to separate files names scene-num-name.go
  * This change went unexpectedly smooth and was done quickly. Go is really awesome for heavy code reshaping! Also putting things in subpackages and naming all magic strings from the beginning helped a  lot.
* Define fyne gray scale colors.
* Not happy with how Options are set in the platform. Refactored it again. Seems to be getting in shape now. Refactoring is still fun!
* Moved all scenes to separate files to see what event loop logic is actually needed to manage a scene. This led to the extraction of some reusabe handlers, such as the global key handler and the options handler. It also allowed moving out the options index logic to a central place.
* Removed GUI platform to focus on termial release
* Moved shared types to `common` package
* Clear game after gameover
