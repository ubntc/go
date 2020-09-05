.PHONY: all test cover cover-check
all: test cover

test:
	go test -race -timeout 3s ./...
	go vet ./...

COVERSH=../scripts/cover.sh
cover:
	$(COVERSH) .

cover-check: cover
	go tool cover -html=.cache/cover.out
