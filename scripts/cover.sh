#/usr/bin/env bash
set -e

cover() { go test ./... -coverpkg "$pkg" -coverprofile=$prefix.out              | tee $prefix.log; }
label() { grep -o 'coverage: [^ \.]*' $prefix.log | sed -e 's/: /-/g'           | tee $prefix.txt; }
link()  { (echo -n "Makefile#"; grep -no "^cover:" Makefile | grep -o '[0-9]*') | tee $prefix.src; }

pkg="$*"
test -n "$pkg" || pkg=./...
pkg=`echo "$pkg" | sed -e 's#\([^ ]*\)#./\1/...#g' -e 's# #,#g'`
prefix=.cache/cover
mkdir -p .cache

echo "covering $pkg"
cover
label=$(label)
link=$(link)
sed -i -e "s/coverage-[0-9]*\(.*\)(.*)$/$label\1($link)/g" README.md
