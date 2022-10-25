package textscenes

const gameover = `

⿴ ⿴ ⿴  ╔═╗╔═╗╔╦╗╔═╗  ╔═╗╦  ╦╔═╗╦═╗  ╦  ⿴ ⿴ ⿴
⿴ ⿴ ⿴  ║ ╦╠═╣║║║║╣   ║ ║╚╗╔╝║╣ ╠╦╝  ║  ⿴ ⿴ ⿴
⿴ ⿴ ⿴  ╚═╝╩ ╩╩ ╩╚═╝  ╚═╝ ╚╝ ╚═╝╩╚═  o  ⿴ ⿴ ⿴

`

var GameOver = Colorize(ColorizeBetween(gameover, "⿴  ", "  ⿴", Bold_Yellow), "⿴", Red)
