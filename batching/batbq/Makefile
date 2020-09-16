.PHONY: all test cover cover-check
all: test cover

test:
	go test -race -timeout 3s ./...
	go vet ./...
	go run _examples/pubsub-to-bq/main.go -dry
	go run _examples/simple/main.go -dry

cover:
	../../scripts/cover.sh ./...

cover-check: cover
	go tool cover -html=.cache/cover.profile
