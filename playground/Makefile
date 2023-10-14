EXAMPLES = bogus corruption display
TESTS    = $(addprefix test-,$(EXAMPLES))

.PHONY: all vars test $(TESTS)

all: vars $(TESTS)

vars:
	# dirs:  $(EXAMPLES)
	# tests: $(TESTS)

cover:
	# Playground code is not covered

OUTPIPE=
test: OUTPIPE=1> /dev/null
test: test-bogus test-display test-corruption

test-bogus:
	go run bogus/bogus.go $(OUTPIPE)
	go run -race bogus/bogus.go -bogus 2>&1 | grep -q "DATA RACE"

test-display:
	go run display/display.go -t 0.5s $(OUTPIPE)

test-corruption:
	go run corruption/corruption.go $(OUTPIPE)

test-foreach:
	go test -bench="ForEach*" ./generics/...
	go test -bench="Loop*"    ./generics/...
