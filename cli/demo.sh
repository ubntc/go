#!/usr/bin/env bash

set -e

cd `dirname $0`  # make sure we are in the right dir

demo=server_demo
proc=server

case "$1" in
    l*)  dir=examples/log;;
    za*) dir=examples/zaplog;;
    ze*) dir=examples/zerolog;;
    *)   dir=examples/zerolog; demo=commands_demo;;
esac

write() {
    for w in $1; do
        echo "$w" | grep -o . | cat | while read char; do
            echo -n "$char"
            sleep 0.02
        done
        echo -n " "; sleep 0.1
    done
    test -z "$2" || sleep "$2"
    echo ""
}
autokill() { 
    (sleep 3.1; kill `pgrep $proc`)& >/dev/null
    "$@"
}
msg() { write "# $1" $2; }
cmd() { write "$*" 1; autokill "$@"; }

server_demo() {
    msg "With a few lines of code Go-cli provides readable colored CLI logs:" 1
    cmd go run $dir/$proc.go -i

    msg "In production we just turn it off:" 1
    cmd go run $dir/$proc.go
    msg "This was a regular run, stopped with CTRL-C." 1

    msg "In interactive mode we can stop with Q, q, CTRL-C, CTRL-D." 1
    cmd go run $dir/$proc.go -i

    write "go run $dir/$proc.go -i" 1
    msg "Watch the clock seconds!" 1
    echo "#                |"
    echo "#                |"
    echo "#                V"
    autokill go run $dir/$proc.go -i

    msg "Wow, that's awesome!" 1
    msg "--------------------------------------"
    msg "Go-cli:" 1
    msg "- easy handling of OS signals" 1
    msg "- easy setup of optional friendly logs"  1
    msg "- easy setup of keyboard commands"  1
    msg "--------------------------------------"
    msg "" 7
    msg "Good bye!"  1
}

commands_demo() {
    msg "Demo app with (h)elp, (s)tatus, and (q)uit commands:"
    go run $dir/$proc.go -i -debug -demo 3h3s3q
    write "" 5
}

$demo