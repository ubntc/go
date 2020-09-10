.PHONY: all test cover cover-check
all: test cover

test:
	go test -race -timeout 3s ./...
	go vet ./...
	go run _examples/pubsub-to-bq/main.go -h 2> /dev/null

cover:
	../../scripts/cover.sh ./...

cover-check: cover
	go tool cover -html=.cache/cover.profile
