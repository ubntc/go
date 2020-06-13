
.PHONY: all demo svg
all:  svg
demo: go-cli.cast go-cli-commands.cast
svg:  go-cli.svg  go-cli-commands.svg

record    = termtosvg record $@
src       = ../demo.sh ../examples/zerolog/server.go
run_demo  = ../demo.sh
theme     = window_frame

%.svg: %.cast ; termtosvg render $^ $@ -t $(theme)

go-cli.cast:          $(src) ; $(record) -g 80x16 -c "$(run_demo) zero"
go-cli-commands.cast: $(src) ; $(record) -g 80x8  -c "$(run_demo)"
