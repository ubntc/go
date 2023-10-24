#!/usr/bin/env bash

set -o errexit

cd "$(dirname "$0")"  # make sure we are in the right dir

demo=server_demo
proc=server

case "$1" in
    l*)  dir=examples/log;;
    za*) dir=examples/zaplog;;
    ze*) dir=examples/zerolog;;
    *)   dir=examples/zerolog; demo=commands_demo;;
esac

gomain="$dir/$proc.go"

write() {
    for w in $1; do
        echo "$w" | grep -o . | cat | while read -r char; do
            echo -n "$char"
            sleep 0.02
        done
        echo -n " "; sleep 0.1
    done
    test -z "$2" || sleep "$2"
    echo ""
}

killserver() {
    sleep 1
    pids="$(pgrep -f "go-build.*/exe/server")"
    # echo "waiting for $gomain PIDs: $pids"
    sleep 2.1
    for pid in $pids; do
        kill "$pid" || true
    done 1> /dev/null 2> /dev/null
}

msg()      { write "# $1" "$2"; }
cmd()      { write "$*" 1; "$@"; }
linenum()  { grep -n "$1" "$2" | cut -d: -f1; }
server()   {
    write "go run $gomain $*" 1
    killserver&
    go run "$gomain" "$@"
}

show_code() {
    start="$(linenum "if \*interactive" "$gomain")"
    end="$(linenum "cli\.WithSigWait" "$gomain")"
    cmd bat -P "$gomain" -r "$start:$end"
}

server_demo() {
    msg "See how easy it is to setup colored logging and interactivity using Go-cli:" 0.5
    show_code
    server -i
    msg "In production we just turn it off:" 1
    server
    msg "This was a regular run, stopped with CTRL-C." 1

    msg "In interactive mode we can stop with Q, q, CTRL-C, CTRL-D." 1
    server -i

    msg "Watch the clock seconds ticking on the bottom line!" 1
    echo "#                |"
    echo "#                |"
    echo "#                V"
    server -i
    msg "Again!" 1
    server -i
    msg "Wow, that's awesome!" 1

    msg "--------------------------------------"
    msg "Go-cli:" 1

    msg "- easy handling of OS signals" 1
    msg "- easy setup of optional friendly logs"  1
    msg "- easy setup of keyboard commands"  1
    msg "- graceful handling of raw terminals"  1
    msg "--------------------------------------"
    msg "" 7
    msg "Good bye!"  1
}

commands_demo() {
    msg "Demo app with (h)elp, (s)tatus, and (q)uit commands:"
    show_code
    go run $gomain -i -debug -demo 3h3s3q
    write "" 5
}

$demo
