# /usr/bin/env bash
set -e

find_tag() {
    git tag --list --no-column --sort=authordate "$1/*" | tail -n 1 | grep -o '[^/]*$'
}

run() { echo "$*"; "$@"; }

for pkg in $*; do
    if tag=`find_tag $pkg`
    then echo "found tag $tag for package $pkg"
    else echo "no tag found for package $pkg"; exit 1
    fi
	# Request the latest tagged version of the package from pkg.go.dev using `go get`.
	# This will trigger an update of the package files in pkg.go.dev and will also add the package
	# to the go.mod file, from which we can remove it using  `go mod tidy`
    echo "fetching $tag for package $pkg"
	run go get -u github.com/ubntc/go/$pkg@$tag
done

run go mod tidy