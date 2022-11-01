.PHONY: run test debug gui install

test:
	go test -race ./...

debug:
	# showing renderings during tests
	go test -v -race ./...

run:
	go run cmd/gotris/gotris.go

gui:
	go run cmd/gotris/gotris.go --platform fyne

install:
	go install cmd/gotris/gotris.go
