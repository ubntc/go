package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/term"
)

type (
	Key      string
	Modifier int
)

type KeyEvent struct {
	Key Key
	Mod Modifier
}

const (
	Shift Modifier = 1 << iota
	Ctrl
	Alt

	Up    Key = "Up"
	Down  Key = "Down"
	Left  Key = "Left"
	Right Key = "Right"
	Quit  Key = "Quit"
	Space Key = "Space"
	Enter Key = "Enter"

	PageUp   Key = "PageUp"
	PageDown Key = "PageDown"
	Home     Key = "Home"
	End      Key = "End"
	Insert   Key = "Insert"
	Delete   Key = "Delete"

	Backspace Key = "Backspace"
	Esc       Key = "Esc"
	Print     Key = "Print"
	Scroll    Key = "Scroll"
	Menu      Key = "Menu"
)

type State struct {
	buffer []byte
	pos    int

	mu sync.RWMutex
}

// Next appends and analyzes input bytes as whole or byte by byte and advances or resets the state.
// Returns a list of all parsed KeyEvents from the known input.
func (s *State) Next(input byte) []KeyEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.buffer = append(s.buffer, input) // Append new input to buffer
	return s.advanceState()
}

func (s *State) advanceState() (events []KeyEvent) {
	defer s.shiftBuffer()

	for s.pos < len(s.buffer) {
		switch s.buffer[s.pos] {
		case '\x1b':
			// handle ESC seq
			seq := s.buffer[s.pos:]

			fmt.Printf("EscCheck:%x: ", seq)

			if len(s.buffer)-s.pos < 3 {
				// found a trailing ESC of length < 3 which may be part of an incomplete ESC seq
				// must stop processing and wait for more input
				// do not advance the position
				return
			}

			// try ESC seq of min 3 bytes to max 5 bytes

			var ev *KeyEvent

			fmt.Printf("EscSeq?:%x: ", seq)

			if ev = ParseKeyEvent3(seq); ev != nil {
				events = append(events, *ev)
				s.pos += 3
				continue
			}
			if ev = ParseKeyEvent4(seq); ev != nil {
				events = append(events, *ev)
				s.pos += 4
				continue
			}
			if ev = ParseKeyEvent6(seq); ev != nil {
				events = append(events, *ev)
				s.pos += 6
				continue
			}

			if HasEscPrefix(seq) {
				return
			}

			fmt.Printf("NoEsc! ")

		case '\xef':
			// Insert Key on Windows Keyboard on Mac in iTerm2

			seq := s.buffer[s.pos:]
			fmt.Printf("Custom len=%d?", len(seq))

			if ev, isPartial := ParseCustomKeyEvent(seq); ev != nil || isPartial {
				if ev != nil {
					events = append(events, *ev)
					s.pos += len(seq)
				}
				if isPartial {
					return
				}
				continue
			}

			fmt.Printf("NoCustom! ")

		}

		fmt.Printf("Char:")
		events = append(events, ParseKeyEvent(s.buffer[s.pos]))
		s.pos++
	}

	return events // Incomplete input, await more bytes
}

func (s *State) shiftBuffer() {
	if s.pos == 0 {
		return
	}
	s.buffer = s.buffer[s.pos:]
	s.pos = 0
}

func HasEscPrefix(input []byte) bool {
	switch {
	case len(input) >= 2 && string(input[:2]) == "\x1b[":
		return true
	default:
		return false
	}
}

func ParseCustomKeyEvent(input []byte) (*KeyEvent, bool) {
	// fmt.Printf("Custom?:%x: ", input)
	switch {
	case len(input) == 3 && string(input) == "\xef\x9d\x86":
		return &KeyEvent{Insert, 0}, false
	case len(input) == 2 && string(input) == "\xef\x9d":
		return nil, true
	case len(input) == 1 && string(input) == "\xef":
		return nil, true
	default:
		return nil, false
	}
}

func ParseKeyEvent3(input []byte) *KeyEvent {
	if len(input) < 3 {
		return nil
	}

	var key Key
	var mod Modifier

	switch string(input[:3]) {

	case "\x1b[A":
		key = Up
	case "\x1b[B":
		key = Down
	case "\x1b[C":
		key = Right
	case "\x1b[D":
		key = Left
	case "\x1b[H":
		key = Home
	case "\x1b[F":
		key = End

	default:
		return nil
	}

	return &KeyEvent{key, mod}
}

func ParseKeyEvent4(input []byte) *KeyEvent {
	if len(input) < 4 {
		return nil
	}

	var key Key
	var mod Modifier

	switch string(input[:4]) {

	case "\x1b[5~":
		key = PageUp
	case "\x1b[6~":
		key = PageDown
	case "\x1b[1~", "\x1b[7~":
		key = Home
	case "\x1b[4~", "\x1b[8~":
		key = End
	case "\x1b[2~": // does not work in iTerm2
		key = Insert
	case "\x1b[3~": // does not work in iTerm2
		key = Delete

	default:
		return nil
	}

	return &KeyEvent{key, mod}
}

func ParseKeyEvent6(input []byte) *KeyEvent {
	if len(input) < 6 {
		return nil
	}

	var key Key
	var mod Modifier

	switch string(input[:6]) {

	case "\x1b[1;2A":
		key = Up
		mod |= Shift
	case "\x1b[1;2B":
		key = Down
		mod |= Shift
	case "\x1b[1;5A":
		key = Up
		mod |= Ctrl
	case "\x1b[1;5B":
		key = Down
		mod |= Ctrl
	case "\x1b[1;5C":
		key = Right
		mod |= Ctrl
	case "\x1b[1;5D":
		key = Left
		mod |= Ctrl
	case "\x1b[1;3A":
		key = Up
		mod |= Alt
	case "\x1b[1;3B":
		key = Down
		mod |= Alt
	case "\x1b[1;3C":
		key = Right
		mod |= Alt
	case "\x1b[1;3D":
		key = Left
		mod |= Alt

	default:
		return nil
	}

	return &KeyEvent{key, mod}
}

func ParseKeyEvent(input byte) KeyEvent {
	// Single-byte character handling

	var key Key

	switch input {

	case 'q', 'Q':
		key = Quit
	case '\x03', '\x04':
		// CTRL+C/D
		key = Quit
	case ' ':
		key = Space
	case '\n', '\x0d':
		key = Enter
	case '\x1b':
		key = Esc
	case '\x7f':
		key = Backspace // ASCII for Backspace
	case '\x10':
		key = Menu
	default:
		key = Key(input)
		// TODO: Do we need explicit support for some chars? Such as: .;-=/\#+-~^°!§$%&()[]{}?´`*'_<>|µ@€'"
	}

	return KeyEvent{key, 0}
}

func captureTerminal() (restore func()) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("MakeRaw: %v", err)
	}
	restore = func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			log.Fatalf("Restore: %v", err)
		}
	}

	return restore
}

func main() {
	restore := captureTerminal()
	defer restore()

	reader := bufio.NewReader(os.Stdin)
	state := State{}
	for {
		// read byte by byte to not block the UI
		b, err := reader.ReadByte()
		if err != nil {
			log.Fatalf("ReadByte: %v", err)
		}

		fmt.Printf("%x:\t", b)

		events := state.Next(b)

		for _, ev := range events {
			fmt.Printf("\tkey:%s", ev.Key)
			if ev.Mod > 0 {
				fmt.Printf("+ mod:%d", ev.Mod)
			}
			fmt.Printf("\n\r")
			if ev.Key == Quit {
				return
			}
		}
	}
}
