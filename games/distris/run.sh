#!/usr/bin/env bash

bin="$(dirname "$0")/bin"

if type -a distris 2> /dev/null; then
    echo "using binary: $(which distris)"
else
    echo "distris not on PATH, adding '$bin' to PATH"
    export PATH="$PATH:$bin"
fi

client() {
    echo "starting client"
    distris -mode client
}

server() {
    echo "starting server"
    distris -mode server
}

async() {
    echo -n "async: "
    "$@"&
    PID=$!
    trap "echo killing server:$PID on SIGTERM; kill $PID" TERM
    trap "echo killing server:$PID on SIGINT;  kill $PID" INT
}

demo() {
    async server
    sleep 0.5
    client
}

if test $# -eq 0
then server || exit 1
else
    for cmd in $*; do
        case $cmd in
            demo)        demo   ;;
            serv*)       async server ;;
            client|play) client ;;
        esac || exit 1
        shift
    done
fi

wait
