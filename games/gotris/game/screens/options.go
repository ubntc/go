package screens

const options = `

        ⿴ ╔═╗╦═╗╔╦╗╦╔═╗╔╗╔╔═╗ ⿴
        ⿴ ║ ║╠═╝ ║ ║║ ║║║║╚═╗ ⿴
        ⿴ ╚═╝╩   ╩ ╩╚═╝╝╚╝╚═╝ ⿴

        Select Rendering Mode:

        MENU_ITEMS

`

var OptionScreen = NewMenuScreen(ColorizeFrame(options, "⿴", Bold_Green, Bold_Red))
