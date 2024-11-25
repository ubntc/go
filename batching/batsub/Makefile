.PHONY: ⚙️
all: ⚙️ test cover

test: ⚙️
	go test -race -timeout 3s ./...
	go vet ./...

cover: ⚙️
	../../scripts/cover.sh ./...

cover-check: ⚙️ cover
	go tool cover -html=.cache/cover.profile
