EXAMPLES = bogus corruption display
TESTS    = $(addprefix test-,$(EXAMPLES))

.PHONY: all vars test $(TESTS)

all: vars $(TESTS)

vars:
	# dirs:  $(EXAMPLES)
	# tests: $(TESTS)

test-bogus:
	go run bogus/bogus.go
	go run -race bogus/bogus.go -bogus 2>&1 | grep "DATA RACE"

test-display:
	go run display/display.go -t 0.5s

test-corruption:
	go run corruption/corruption.go
