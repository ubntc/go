// The code in this file was generated in 2022-11 using https://beta.openai.com/playground
// using the text-davinci-002 model.
// Please do not copy consider this any form of production ready code.
// It is just a demo to show how the OpenAI models can generate runnable code
// and actually teach you how to use some libs or solve common tasks.
//
// Do not copy this code into your learning code base for your AI!
// It is really rough and incomplete.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/eiannone/keyboard"
)

func main() { davinci_002_game_05() }

//nolint
func davinci_002_game_01() {
	c := exec.Command("stty", "-icanon", "min", "1")
	c.Stdin = os.Stdin
	c.SysProcAttr = &syscall.SysProcAttr{}
	c.Run()
	defer c.Process.Kill()

	var b []byte = make([]byte, 1)
	for {
		_, _ = os.Stdin.Read(b)
		fmt.Println("you pressed:", b, string(b))
	}
}

//nolint
func davinci_002_game_02() {
	fmt.Println("Welcome to the game! Use the arrow keys to move the character.")

	// Create the game character
	type character struct {
		x int
		y int
	}
	var player character

	// Read in key presses from the terminal
	for {
		var input string
		fmt.Scanln(&input)

		// Map the key presses to movement actions
		switch input {
		case "up":
			player.y--
		case "down":
			player.y++
		case "left":
			player.x--
		case "right":
			player.x++
		case "q":
			return
		default:
			fmt.Println("Please use the arrow keys to move the character.")
		}

		// Print the game character's new position
		fmt.Printf("Character position: (%d, %d)\n", player.x, player.y)
	}
}

//nolint
func davinci_002_game_03() {
	c := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1")
	c.Stdin = os.Stdin
	c.Run()

	defer func() {
		c := exec.Command("stty", "-F", "/dev/tty", "echo")
		c.Stdin = os.Stdin
		c.Run()
	}()

	getch := func() byte {
		var b [1]byte
		os.Stdin.Read(b[:])
		return b[0]
	}

	width, height, err := terminalSize()
	if err != nil {
		fmt.Println(err)
		return
	}

	x, y := width/2, height/2

	for {
		clearScreen()
		fmt.Printf("\033[%d;%dH*", y, x)

		switch getch() {
		case '\033': // if the first value is esc
			getch() // skip the [
			switch getch() {
			case 'A':
				y--
			case 'B':
				y++
			case 'C':
				x++
			case 'D':
				x--
			}
		case 'q':
			return
		}
	}
}

//nolint
func davinci_002_game_04() {
	// Read input from the terminal.
	input := os.Stdin

	// Loop forever.
	for {
		// Read a single byte from the input.
		b := make([]byte, 1)
		var err error
		_, err = input.Read(b)
		if err != nil {
			panic(err)
		}

		// Process the byte.
		switch b[0] {
		case 27: // ESC
			fmt.Println("ESC")
		case 91: // [
			fmt.Println("[")
		case 68: // Left arrow.
			fmt.Println("Left arrow")
		case 67: // Right arrow.
			fmt.Println("Right arrow")
		case 65: // Up arrow.
			fmt.Println("Up arrow")
		case 66: // Down arrow.
			fmt.Println("Down arrow")
		default:
			fmt.Println("Other")
		}
	}
}

//nolint
func davinci_002_game_05() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyArrowUp:
			fmt.Println("Up")
		case keyboard.KeyArrowDown:
			fmt.Println("Down")
		case keyboard.KeyArrowLeft:
			fmt.Println("Left")
		case keyboard.KeyArrowRight:
			fmt.Println("Right")
		}

		if char == 'q' {
			break
		}
	}
}

//nolint
func clearScreen() {
	cmd := exec.Command("clear") // Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//nolint
func terminalSize() (width, height int, err error) {
	ws := &Winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		os.Stdin.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		err = errno
		return
	}
	width = int(ws.Col)
	height = int(ws.Row)
	return
}

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}
