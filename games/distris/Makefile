.PHONY: build client server demo

SOURCES := $(wildcard **/*.go) main.go go.mod go.sum
DISTRIS := bin/distris

$(DISTRIS): $(SOURCES)
	# building binary: $@
	# from sources: $^
	weaver generate ./...
	go mod download
	go build -o $@ main.go

client demo server: $(DISTRIS)
	./run.sh $@ || true
