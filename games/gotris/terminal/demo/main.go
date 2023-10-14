package main

import (
	"context"
	"os"
	"time"

	"github.com/ubntc/go/games/gotris/terminal"
)

// func scan() {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		fmt.Println("key:", scanner.Text())
// 	}
// }

func main() {
	t := terminal.NewTerminal(os.Stdout)
	input, restore, err := t.CaptureInput(context.Background())
	if err != nil {
		panic(err)
	}
	defer restore()

	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	clearCountAvg := 0.0
	clearCount := 0
	timeWindowStart := time.Now()

	pause := true
	debug := os.Getenv("DEBUG") != ""

	echo := t.Println
	help := func() {
		echo("-------------------------------------------------")
		echo("Press 'c', 'x', or 'o' to test the clear methods.")
		echo("Hold the key to test the clear speed.            ")
		echo("                                                 ")
		echo("Press 'h' to show this help (and pause).         ")
		echo("Press 'p' to start and pause.                    ")
		echo("-------------------------------------------------")
	}

	showResults := func() {
		dt := time.Since(timeWindowStart)
		if dt > time.Second {
			timeWindowStart = time.Now()
			clearCountAvg = clearCountAvg*0.5 + float64(clearCount*int(time.Second))/float64(dt)
			clearCount = 0
		}
		echo("Clears in current Second:", clearCount)
		echo("Clears per Second:       ", clearCountAvg)
	}

	help()

	i := 0
	for {
		select {
		case in, more := <-input:
			if !more {
				return
			}
			if debug {
				echo("key: ", in.Key(), ", flags: ", in.Flags(), ", rune: ", in.Text(), ", mov: ", in.IsMovement(), ", alt: ", in.IsAlt())
			}
			switch in.Rune() {
			case 'q':
				echo("Quit")
				return
			case 'c':
				t.Clear()
				echo("Cleared")
			case 'o':
				_ = t.Overpaint()
				echo("Overpainted")
			case 'x':
				t.RunClearCommand()
				echo("Executed Clear Command")
			case 'h':
				t.Clear()
				pause = true
				help()
				continue
			case 'p':
				pause = !pause
				if pause {
					echo("PAUSE")
				}
				continue
			default:
				continue
			}
			clearCount++
			showResults()

		case <-ticker.C:
			if !pause {
				echo(i, time.Now())
			}
		}
		i += 1
	}
}
