package textscenes

var welcome = Colorize(`

        ⿴  ╔═╗╔═╗╔╦╗╦═╗╦╔═╗  ⿴
        ⿴  ║ ╦║ ║ ║ ╠╦╝║╚═╗  ⿴
        ⿴  ╚═╝╚═╝ ╩ ╩╚═╩╚═╝  ⿴

        A fancy block drop game
        for your system terminal!

          MENU_ITEMS

╭[i]────────────────────────────────────────────────╮
│ ©2022 @ubunatic. Made by with ♡ in Saxony.        │
│ Code at https://github.com/ubntc/go/games/gotris. │
╰───────────────────────────────────────────────────╯
`, "♡", Bold_Red)

var WelcomeScreen = NewMenuScreen(ColorizeFrame(welcome, "⿴", Bold_Green, Bold_Red))
