# This Makefile implements basic build and test actions for the ubntc/go projects.

.PHONY: help debug test cover make tag refresh

PACKAGES = errchan cli playground batching/batsub batching/batbq
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
refresh: $(addsuffix -refresh, $(TARGETS))

# Generic package testing and coverage
%-test:    %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/test.sh $<
%-cover:   %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/cover.sh $<

# Generic build targets and version management
%-make:    %; $(MAKE) -C $<
%-tag:     %; # TODO: use bumpversion tool + push it
%-refresh: %; scripts/version.sh $<
