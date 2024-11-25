#!/usr/bin/env bash
set -e

coverage() {
    if test -e "$prefix.profile"; then
        go tool cover -func "$prefix.profile" -o "$prefix.txt"
        tail -n1 "$prefix.txt" | grep -o '[0-9]\+' | head -n1
    fi
}

link() {
    echo -n "Makefile#"
    grep -no "^cover:" Makefile | grep -o '[0-9]\+'
}

update_readme() {
    set -e
    if test -e README.md; then
        if test -e "Makefile"
        then link=$(link)
        else link=
        fi
        label=$(coverage)
        test -n "$label" || label="0"
        echo "setting label=$label and link=$link in $PWD/README.md"
        sed -i -e "s/coverage-[0-9\.]*%25\(.*\)(.*)$/coverage-$label%25\1($link)/g" README.md
        echo -n "resulting badge: "
        if ! grep -e "coverage-[0-9]\+%25" README.md
        then echo "no badge found"
        fi
    fi
}

test -n "$prefix" || prefix=.cache/cover
# update_readme #TODO: fix genertation on MacOS
