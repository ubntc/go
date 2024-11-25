#!/usr/bin/env bash
set -e

# find_tag finds the latest Go tag for a given subpackage's git tag.
find_tag() {
    git tag --list --no-column --sort=authordate "$1/*" | tail -n 1 | grep -o '[^/]*$'
}

# run prints a command before running it, similar to how `make` echos commands.
run() { echo "$*"; "$@"; }

# touch_tag queries pkg.go.dev for the latest tag of a given `pkg`
# to trigger an update of the published package.
touch_tag(){
    pkg=$1
    if tag=$(find_tag "$pkg")
    then echo "found tag $tag for package $pkg"
    else echo "no tag found for package $pkg"; exit 1
    fi
    if test -e "$pkg/go.mod"
    then echo "found $pkg/go.mod"
    else echo "$pkg/go.mod not found, $pkg is not a standalone module"; exit 1
    fi
    # Request the latest tagged version of the package from pkg.go.dev using `go get`.
    # This will trigger an update of the package files in pkg.go.dev and will also add the package
    # to the go.mod file, from which we can remove it using  `go mod tidy`
    echo "fetching $tag for package $pkg"
    run go get "github.com/ubntc/go/$pkg@$tag"
}

trap "run go mod tidy" EXIT              # ensure cleanup
for pkg in "$@"; do touch_tag "$pkg"; done   # touch all packages
