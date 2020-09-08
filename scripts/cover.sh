#/usr/bin/env bash
set -e

run(){ echo "$*" 1>&2; "$@"; }

cover() {
    run go test $GOTEST_ARGS -coverpkg "$pkg" -coverprofile=$prefix.profile ./...
    prefix=$prefix $here/badges.sh
}

coverpkg="."  # subpackages to be covered
packages=""  # standalone packages to be covered separately
COVERSH=`readlink -f "$0"`
here=`dirname $COVERSH`
prefix=.cache/cover
mkdir -p .cache

for p in $*; do
    if test -e "$p/go.mod" && test "$p" != .
    then packages="$packages $p"
    else coverpkg="$coverpkg,./$p/..."
    fi
done

for p in $packages; do
    echo "covering standalone package $p"
    if test -e Makefile && grep -q "^cover:" Makefile
    then make -C $p cover COVERSH=$COVERSH
    else cd $p && pkg=./... cover
    fi
done

if test -n "$coverpkg"; then
    echo "covering packages $coverpkg"
    test -z "$GOTEST_ARGS" || echo "using GOTEST_ARGS=$GOTEST_ARGS as additional cover arguments"
    pkg=$coverpkg cover
fi
