#!/usr/bin/env bash

here=$(dirname "$0")
rsync -a "$here/corruption.go pi@raspberrypi:/tmp/"
ssh -C pi@raspberrypi go run /tmp/corruption.go "$@"
