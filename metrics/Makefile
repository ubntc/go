.PHONY: all cover test
all: test cover
cover:
	../scripts/cover.sh ./...

test:
	go test -race ./...
