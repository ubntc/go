# (!) Generated File (!)
# Consider making changes in the main Makefile instead.
#
# This is a generic file to help developing Go projects.
# It provides file and command targets for common Go dev tasks.

# Sources and Binaries
# --------------------
project := $(notdir $(CURDIR))
main = $(project).go
binary = bin/$(project)
sources = $(shell find . -name '*.go' -o -name 'go.*')

$(binary): $(sources); go build -o $@ $(main)

# Common Targets
# --------------
.PHONY: build test debug install run vars usage all
all:   build
build: $(binary)

test:    build; go test -race ./...
debug:   build; go test -v -race ./...
install: build; go install .
run:     build; $(binary)

vars:
	# project  $(project)
	# main:    $(main)
	# binary:  $(binary)
	# sources: $(sources)

usage:
	# Usage: make TARGET [-n|-B]
	#
	# Targets:
	#     build     builds the binary: $(binary)
	#     test      runs go test
	#     debug     runs go test -v
	#     run       builds and runs the binary $(binary)
	#     install   runs go install
	#     vars      shows the make vars
	#     usage     shows this info
	#
	# Flags:
	#     -n        dry runs, just show commands that would run
	#     -B        force build even if sources have not changed
	#
	#     Also see `make -h` for more common flags
