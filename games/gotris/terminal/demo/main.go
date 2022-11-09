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
	t := terminal.New(os.Stdout)
	input, restore, err := t.CaptureInput(context.Background())
	if err != nil {
		panic(err)
	}
	defer restore()

	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	timeFirstClear := time.Now()
	clearCount := 0
	pause := true
	echo := t.Println
	debug := os.Getenv("DEBUG") != ""

	echo("-------------------------------------------------")
	echo("Press 'c', 'x', or 'o' to test the clear methods.")
	echo("Hold the key to test the clear speed.            ")
	echo("Press 'p' to start.                              ")
	echo("-------------------------------------------------")

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
			case 'p':
				pause = !pause
				if pause {
					echo("PAUSE")
				}
				continue
			default:
				continue
			}
			// print clear speed results
			dt := time.Since(timeFirstClear)
			clearCount += 1
			cps := clearCount * int(time.Second) / int(dt)
			echo("Clears/Second:", cps)
		case <-ticker.C:
			if !pause {
				echo(i, time.Now())
			}
		}
		i += 1
	}
}
