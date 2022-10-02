.PHONY: run test

test:
	go test -race ./...

debug:
	# showing renderings during tests
	go test -v -race ./...

run:
	go run cmd/gotris/gotris.go

install:
	go install cmd/gotris/gotris.go
