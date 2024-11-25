package config

type Config struct {
	ShowClock   bool // ShowClock enables the default ascii/unicode clock in the last terminal line.
	WithQuit    bool // WithQuit adds the default quit commands and enables user input.
	PrependCR   bool
	MakeTermRaw bool
}

func Default(interactive bool) Config {
	if interactive {
		return Interactive()
	}
	return Server()
}

func Server() Config {
	return Config{
		ShowClock:   false,
		WithQuit:    false,
		PrependCR:   false,
		MakeTermRaw: false,
	}
}

func Interactive() Config {
	return Config{
		ShowClock:   true,
		WithQuit:    true,
		PrependCR:   true,
		MakeTermRaw: true,
	}
}
