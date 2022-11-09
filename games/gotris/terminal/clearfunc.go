package terminal

import (
	"os"
	"os/exec"
	"runtime"
)

// clearfunc source code is based on
// https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go/53673326#53673326

var clear map[string]func(*os.File) // create a map for storing clear funcs

func init() {
	clear = make(map[string]func(*os.File))
	clear["linux"] = func(fd *os.File) {
		cmd := exec.Command("tput", "clear")
		cmd.Stdout = fd
		_ = cmd.Run()
	}

	clear["darwin"] = clear["linux"]

	clear["windows"] = func(fd *os.File) {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = fd
		_ = cmd.Run()
	}
}

func callClearFunc(fd *os.File) {
	// works in iTerm2 with some flickering and perfectly in VSCode terminal
	if clearFunc, ok := clear[runtime.GOOS]; ok {
		clearFunc(fd)
		return
	}
	panic("Platform is unsupported! Can't clear Terminal screen.")
}
