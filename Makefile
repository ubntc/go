# This Makefile implements basic build and test actions for the projects.
# It tests and covers Go code per project, using the projects own `Makefile`,
# or via regular `go test` commands, or using the `./scripts`.

.PHONY: all help test cover

PACKAGES = errchan cli playground
TARGETS  = cli

_TESTS    = $(addsuffix -test,  $(PACKAGES))
_COVERS   = $(addsuffix -cover, $(PACKAGES))
_MAKES    = $(addsuffix -make,  $(TARGETS))

help:
	# usage: make TARGET PACKAGES="dir1 dir2 ..."
	# targets:
	#
	#   test      run tests for PACKAGES=$(PACKAGES)
	#   cover     check and update test coverage for PACKAGES=$(PACKAGES)
	#   all       run `make` in subdirs for TARGETS=$(TARGETS)

debug: GOTEST_ARGS=-v
debug: test
test: $(_TESTS)
cover: $(_COVERS)
all: $(_MAKES)

# Generic go test to test all subpackages.
%-test:  %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/test.sh $<
%-cover: %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/cover.sh $<
%-make:  %; $(MAKE) -C $<
