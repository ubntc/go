.PHONY: all clean test debug resources done cover cover-check
all: test cover resources

clean: 
	rm resources/go-cli*.cast resources/go-cli*.svg .cache/cover*

test:
	go test -race -timeout 3s ./...
	go vet ./...

debug:
	go test -race -failfast -v -timeout 3s ./...

done:
	# marks demo files as new to avoid regeneration
	touch resources/go-cli-*

resources:
	$(MAKE) -C $@

COVERSH=echo
cover:
	$(COVERSH) cli loggers

cover-check: cover
	go tool cover -html=.cache/cover.out
