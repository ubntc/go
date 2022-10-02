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

Run it with:

```
go run cmd/main.go
```


## Implementation History

Features were implemented in the following order.

 1. Basic game enties and game loop
 2. Tile movement and collision
 3. Basic terminal rendering
 4. User input (via stdin)
 5. Tile rotation (which triggerd quite some refactoring)
 6. Smart tile drawing using a short string to define how to draw a tile (instead of slices of numbers)
 7. Drawing specs for orientations with a well-defined center (it is more complicated than you may think)
 8. Smart rotation (move tile left or right to resolve blocked rotation)
 9. Scoring (only then the game was playable and fun 🥳🎊🎉)
10. Next tile preview
11. Rendering components (for preview and score box)
12. Frames around view boxes
13. Better speed adjustment
14. Separation of concers: dissolve code of big main.go into the packages
15. Make game package independent of the renderer and terminal (added "platform" notion)

## Missing Features

1. Sound
2. Official scoring rules
