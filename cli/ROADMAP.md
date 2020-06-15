# Go-cli Roadmap

## Musthaves
* [x] Define scope of project
* [x] Automatic race condition tests
* [x] Automatic tests
* [x] Test coverage check
* [x] Automatically recordable ASCII demo of the output

### Logging
* [x] Define log + time format
* [x] Setup Stdlog
* [x] Setup Zerolog
* [ ] Setup Zaplog
* [ ] Auto-Wrap Stdlog
* [x] Examples + Readme

### Signals 
* [x] Catch + terminate on `os.Signal` and stop app by closing the context
* [x] Ensure graceful shutdown after closing context
* [ ] Test external signal handling (TERM, KILL)
* [ ] Test user interruption handling (SIGINT, ^C)

### Input + Commands
* [x] Allow user input
* [x] Allow single key input
* [x] Ensure terminal is restored on termination
* [x] Bind default quit keys Q,q,^C,^D
* [x] Custom commands + key binds
* [x] Run script (sequence of keys)
* [ ] Ensure terminal is restored on various panic scenarios

### Compatibility
* [x] Runs on Desktop Linux
* [x] Runs on RaspberryPi
* [ ] Runs in Docker
* [ ] Runs on Mac (iterm)
* [ ] Runs on Windows (cmd)
* [ ] Runs in GCP

## Nicetohaves
* [x] Run scripts by chaining commands
* [ ] Generic logger wrapping (e.g., use zero with wrapped std + zap)
* [ ] Define how to interact with generic `Logger` 
* [ ] Generic Setup for any `Logger`
* [ ] Setup for logrus
* [ ] More supported loggers
* [ ] Detect unicode support and fallback to ASCII clock
* [ ] Use go-termios directly to improve restoring terminal state
* [ ] Detect broken terminal state and repair
* [ ] Runs on all Go-supported platforms
