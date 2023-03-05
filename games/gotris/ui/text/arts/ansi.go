package arts

const (
	Reset = "\x1b[0m"

	Framed = "\x1b[51m" // does not work in iTerm and VSCode term
	Bold   = "\x1b[1m"
	Invert = "\x1b[7m"

	FG_Orange = "\x1b[91m"
	FG_Red___ = "\x1b[31m"
	FG_Green_ = "\x1b[32m"
	FG_Yellow = "\x1b[33m"
	FG_Blue__ = "\x1b[34m"
	FG_Magent = "\x1b[35m"
	FG_Purple = "\x1b[94m"
	FG_Pink__ = "\x1b[95m"
	FG_Cyan__ = "\x1b[36m"

	FG_Black_ = "\x1b[30m"
	FG_D_Gray = "\x1b[38;5;234m"
	FG_Gray__ = "\x1b[90m"
	FG_White_ = "\x1b[97m" // bright white
	FG_L_Gray = "\x1b[37m" // dull white

	BG_Orange = "\x1b[101m"
	BG_Red___ = "\x1b[41m"
	BG_Green_ = "\x1b[42m"
	BG_Yellow = "\x1b[43m"
	BG_Blue__ = "\x1b[44m"
	BG_Magent = "\x1b[45m"
	BG_Purple = "\x1b[104m"
	BG_Pink__ = "\x1b[105m"
	BG_Cyan__ = "\x1b[46m"

	BG_Black_ = "\x1b[40m"
	BG_D_Gray = "\x1b[48;5;233m"
	BG_Gray__ = "\x1b[100m"
	BG_White_ = "\x1b[107m" // bright white
	BG_L_Gray = "\x1b[47m"  // dull white
)
