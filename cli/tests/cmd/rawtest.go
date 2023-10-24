// ioctldevice tests cannnot be run via `go test`.
// We need to run this CLI tool to check the funtionality.
package main

import (
	"flag"
	"log"

	"github.com/ubntc/go/cli/cli"
	"golang.org/x/sys/unix"
)

// TODO: Do we need a fallback mode for broken terminals?
// If so, we need to collect common termios states.
// X1 Carbon, Arch Linux, Interactive Mode in Terminal:
var (
	FallbackStateTilixTerminal = unix.Termios{
		Iflag: 17664, Oflag: 5, Cflag: 191, Lflag: 35387,
		Cc: [20]uint8{3, 28, 127, 21, 4, 0, 1, 0, 17, 19, 26, 0, 18, 15, 23, 22, 0, 0, 0, 0},
	}

	FallbackStateVSCodeTerminal = unix.Termios{
		Iflag: 26624, Oflag: 4, Cflag: 1215, Lflag: 2618,
		Cc: [20]uint8{3, 28, 127, 21, 4, 0, 1, 0, 17, 19, 26, 255, 18, 15, 23, 22, 255, 0, 0, 0},
	}

	// MacOS VSCode
	FallbackStateMacOSVSCodeTerminal = unix.Termios{
		Iflag: 27394, Oflag: 3, Cflag: 19200, Lflag: 536872399,
		Cc:     [20]uint8{4, 255, 255, 127, 23, 21, 18, 0, 3, 28, 26, 25, 17, 19, 22, 15, 1, 0, 20, 0},
		Ispeed: 38400, Ospeed: 38400,
	}
)

func main() {
	interative := flag.Bool("i", false, "interactive mode for ioctl test")
	flag.Parse()
	term := cli.GetTerm()
	log.SetOutput(term)
	term.SetDebug(true)
	term.SetVerbose(true)
	restore, err := cli.ClaimTerminal()
	if restore != nil {
		defer restore() // nolint
	}
	defer term.Println("rawtest done")
	log.Println("rawtest: ResoreFunc:", restore)
	log.Println("rawtest: error:", err)
	switch *interative {
	case true:
		if restore == nil {
			log.Fatalln("rawtest: restore func is nil")
		}
		if err != nil {
			log.Fatalln("rawtest: failed to claim terminal")
		}
	case false:
		// assuming to run in non-interactive mode. ClaimTerminal must fail now!
		if restore != nil {
			log.Println("rawtest: WARNING! ResoreFunc is not nil")
		}
		if err == nil {
			log.Fatalln("rawtest: ClaimTerminal error should not be nil")
		}
	}
}
