# This Makefile implements basic build and test actions for the projects.
# It tests and covers Go code per project, using the projects own `Makefile`,
# or via regular `go test` commands, or using the `./scripts`.

.PHONY: help debug test cover make tag publish

PACKAGES = errchan cli playground
TARGETS  = cli

help:
	# usage: make TARGET PACKAGES="dir1 dir2 ..."
	# targets:
	#
	#   test      run tests for PACKAGES=$(PACKAGES)
	#   cover     check and update test coverage for PACKAGES=$(PACKAGES)
	#   make      run `make` in subdirs for TARGETS=$(TARGETS)

debug: GOTEST_ARGS=-v
debug: test
test:    $(addsuffix -test,    $(PACKAGES))
cover:   $(addsuffix -cover,   $(PACKAGES))
make:    $(addsuffix -make,    $(TARGETS))
tag:     $(addsuffix -tag,     $(TARGETS))
publish: $(addsuffix -publish, $(TARGETS))

# Generic go test to test all subpackages.
%-test:    %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/test.sh $<
%-cover:   %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/cover.sh $<
%-make:    %; $(MAKE) -C $<

# Version management
_RAWTAG = $(shell git tag --list --no-column --sort=authordate '$</v*' | tail -n 1 | grep '^$</v.*$$')
_TAG    = $(shell echo '$(_RAWTAG)' | cut -d '/' -f 2)

%-tag:     %; # TODO: use bumpversion lib, current TAG=$(_TAG)
%-publish: %; echo curl https://sum.golang.org/lookup/github.com/go/$<@$(_TAG)
