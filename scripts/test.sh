#/usr/bin/env bash
set -e

testpkg() {
    set -e
    pkg="$1"
    echo -n "testing '$pkg', "
    if test -e "$pkg/Makefile"; then
        echo "found '$pkg/Makefile'"
        run make -C "$pkg" test
    else
        if test -e "$pkg/go.mod"
        then tests="./...";       echo "found go.mod for package '$pkg', entering dir '$pkg'"; cd "$pkg"
        else tests="./$pkg/...";  echo "found subpackage '$pkg', running subpackage tests";
        fi

        mkdir -p .cache
        prefix=.cache/"$pkg"
        test -z "$GOTEST_ARGS" || echo "using GOTEST_ARGS=$GOTEST_ARGS as additional go test arguments"
        run go test -race $GOTEST_ARGS $tests 2>&1 | tee $prefix.log
        run go vet $tests
    fi

    if test -e "$pkg/go.mod" && ! test -f "$pkg/LICENSE"; then
        echo "stand-alone package $pkg does not have hard copied LICENSE file"
        return 1
    fi
}

run() { echo "$*"; "$@"; }

for pkg in $*; do testpkg "$pkg"; done
