.PHONY: all test cover cover-check example
all: test cover example

test:
	go test -race -timeout 3s ./...
	go vet ./...
	go run _examples/ps2bq/main.go -dry
	go run _examples/simple/main.go -dry

cover:
	../../scripts/cover.sh ./...

cover-check: cover
	go tool cover -html=.cache/cover.profile

example:
	go run ../../scripts/code_include.go README.md
