# This Makefile implements basic build and test actions for the ubntc/go projects.

.PHONY: ⚙️

MODULES  = $(shell find . -name go.mod | xargs dirname)
PACKAGES = internal/errchan cli playground batching/batsub batching/batbq metrics
TARGETS  = cli batching/batsub batching/batbq metrics

help: ⚙️
	# usage: make TARGET PACKAGES="dir1 dir2 ..."
	# targets:
	#
	#   test      run tests for PACKAGES=$(PACKAGES)
	#   cover     check and update test coverage for PACKAGES=$(PACKAGES)
	#   tidy      run go mod tidy for MODULES=$(MODULES)
	#   make      run `make` in subdirs for TARGETS=$(TARGETS)

debug: GOTEST_ARGS=-v
debug: ⚙️ test

test:    ⚙️ $(addsuffix -test,    $(PACKAGES))
cover:   ⚙️ $(addsuffix -cover,   $(PACKAGES))
make:    ⚙️ $(addsuffix -make,    $(TARGETS))
tag:     ⚙️ $(addsuffix -tag,     $(TARGETS))
refresh: ⚙️ $(addsuffix -refresh, $(TARGETS))

tidy: ⚙️ $(addsuffix -tidy, $(MODULES))
	go work sync

# Generic package testing and coverage
%-test:    %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/test.sh $*
%-cover:   %; GOTEST_ARGS="$(GOTEST_ARGS)" scripts/cover.sh $*
%-tidy:    %; cd $* && go mod tidy

# Generic build targets and version management
%-make:    %; $(MAKE) -C $*
%-tag:     %; # TODO: use bumpversion tool + push it
%-refresh: %; scripts/version.sh $*
