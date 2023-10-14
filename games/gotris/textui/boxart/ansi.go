package boxart

const (
	Reset = "\x1b[0m"

	Framed = "\x1b[51m" // does not work in iTerm and VSCode term
	Bold   = "\x1b[1m"
	Invert = "\x1b[7m"

	// To preview colors here, add "**/*.go" to the "colorize.include" setting in VSCode
	// for the kamikillerto.vscode-colorize plugin

	Ora = "\x1b[91m" // #ff4500
	Red = "\x1b[31m" // #ff0000
	Grn = "\x1b[32m" // #008000
	Yel = "\x1b[33m" // #ffff00
	Blu = "\x1b[34m" // #0000ff
	Mag = "\x1b[35m" // #ff00ff
	Pur = "\x1b[94m" // #800080
	Pnk = "\x1b[95m" // #ff69b4
	Cyn = "\x1b[36m" // #00ffff

	Gry = "\x1b[90m" // #808080
	Blk = "\x1b[30m" // #000000
	Wht = "\x1b[97m" // #FFFFFF

	DkGry = "\x1b[38;5;234m" // #1c1c1c
	LtGry = "\x1b[37m"       // #d3d3d3

	BgOra = "\x1b[101m" // #ff4500
	BgRed = "\x1b[41m"  // #ff0000
	BgGrn = "\x1b[42m"  // #008000
	BgYel = "\x1b[43m"  // #ffff00
	BgBlu = "\x1b[44m"  // #0000ff
	BgMag = "\x1b[45m"  // #ff00ff
	BgPur = "\x1b[104m" // #800080
	BgPnk = "\x1b[105m" // #ff69b4
	BgCyn = "\x1b[46m"  // #00ffff

	BgBlk = "\x1b[40m"  // #000000
	BgWht = "\x1b[107m" // #ffffff
	BgGry = "\x1b[100m" // #808080

	BgDkGry = "\x1b[48;5;233m" // #1c1c1c
	BgLtGry = "\x1b[47m"       // #d3d3d3

)
