# Go-cli Roadmap

## Musthaves
* [x] Define scope of project
* [x] Automatic race condition tests
* [x] Automatic tests
* [x] Test coverage check
* [x] Automatically recordable ASCII demo of the output

### Logging
* [x] Define log + time format
* [x] Setup Std Log
* [x] Setup Zerolog
* [ ] Setup Zaplog
* [x] Examples + Readme
* [ ] Define how to interact with std `Logger` interface and create "Anylogger" Setup

### Signals 
* [x] Catch + terminate on `os.Signal` and stop app via closing context
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
* [ ] Runs on Windows
* [ ] Runs on Mac
* [ ] Runs in GCP 

## Nicetohaves
* [ ] Setup logrus
* [ ] Setup logger 1
* [ ] Setup logger 2
* [ ] Setup logger 3
* [ ] Detect unicode support and fallback to ASCII clock