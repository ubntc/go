.PHONY: ⚙️

EXAMPLES = bogus corruption display foreach generics diff compat
TESTS    = $(addprefix test-,$(EXAMPLES))

all: vars $(TESTS)

vars: ⚙️
	# dirs:  $(EXAMPLES)
	# tests: $(TESTS)

cover: ⚙️
	# Playground code is not covered

OUTPIPE=
test: OUTPIPE=1> /dev/null
test: ⚙️ $(TESTS)

test-%: ⚙️
	go test -v ./$* $(OUTPIPE)

test-bogus: ⚙️
	go run bogus/bogus.go $(OUTPIPE)
	go run -race bogus/bogus.go -bogus 2>&1 | grep -q "DATA RACE"

test-display: ⚙️
	go run display/display.go -t 0.5s $(OUTPIPE)

test-corruption: ⚙️
	go run corruption/corruption.go $(OUTPIPE)

test-foreach: ⚙️
	go test -bench="All" ./generics/...
