.PHONY: all clean test debug resources done cover cover-check
all: test cover resources

clean:
	rm resources/go-cli*.cast resources/go-cli*.svg .cache/cover*

test:
	go test -race -timeout 3s ./...
	go vet ./...
	# INTERACTIVE MODE TEST
	go run -race tests/cmd/rawtest.go -i
	# NON-INTERACTIVE MODE TEST
	echo 1 | go run -race tests/cmd/rawtest.go

debug:
	go test -race -failfast -v -timeout 3s ./...

done:
	# marks demo files as new to avoid regeneration
	touch resources/go-cli-*

resources:
	$(MAKE) -C $@

cover:
	../scripts/cover.sh cli loggers

cover-check: cover
	go tool cover -html=.cache/cover.profile
