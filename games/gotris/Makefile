.PHONY: run test debug install

test:
	go test -race ./...

debug:
	# showing renderings during tests
	go test -v -race ./...

run:
	go run gotris.go

install:
	go install gotris.go
