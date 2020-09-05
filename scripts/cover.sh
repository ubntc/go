#/usr/bin/env bash
set -e

cover() { go test ./... $GOTEST_ARGS -coverpkg "$pkg" -coverprofile=$prefix.out   | tee $prefix.log; }
label() { grep -o 'coverage: [^ \.]*' $prefix.log | sed -e 's/: /-/g' | tail -n 1 | tee $prefix.txt; }
link()  { (echo -n "Makefile#"; grep -no "^cover:" Makefile | grep -o '[0-9]*')   | tee $prefix.src; }

update_readme() {
    set -e
    if test -e README.md; then
        if test -e "Makefile"
        then link=$(link)
        else link=
        fi
        label=$(label)
        echo "setting label=$label and link=$link in README.md"
        sed -i -e "s/coverage-[0-9]*\(.*\)(.*)$/$label\1($link)/g" README.md
    fi
}

coverpkg=""  # subpackages to be covered
packages=""  # standalone packages to be covered separately
script=`readlink -f "$0"`
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
    if test -e Makefile
    then make -C $p cover COVERSH=$script
    else cd $p && pkg=./... prefix=.cache/"$p-cover" cover
    fi
done

if test -n "$coverpkg"; then
    echo "covering packages $coverpkg"
    test -z "$GOTEST_ARGS" || echo "using GOTEST_ARGS=$GOTEST_ARGS as additional cover arguments"
    pkg=$coverpkg cover
    update_readme
fi
