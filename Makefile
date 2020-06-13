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

test: $(_TESTS)
cover: $(_COVERS)
all: $(_MAKES)

# Generic go test to test all subpackages.
%-test: %;
	go test -race $(GOCOVER_ARGS) ./$</...
	go vet  ./$</...

COVERSH_CMD = $(MAKE) -C $< cover COVERSH=$(CURDIR)/scripts/cover.sh
# The above COVERSH_CMD is used for projects that have their own Makefile.
# Projects without Makefile use `go test -cover` to compute test coverage.
%-cover: GOCOVER_ARGS = -cover -coverprofile=.cache/$<.out 
%-cover: %
	if test -e $</Makefile; then $(COVERSH_CMD); else $(MAKE) $<-test; fi

# Runs the default make target of a project.
%-make: %; $(MAKE) -C $<
