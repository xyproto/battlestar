#!/bin/sh
battlestar=../battlestar
if [ ! -e $battlestar ]; then
  echo "Could not find $battlestar!"
  exit 1
fi
bits=`getconf LONG_BIT`
osx=$([[ `uname -s` = Darwin ]] && echo true || echo false)
for f in *.bts; do
  n=${f/.bts}
  "$battlestar" -bits="$bits" -osx="$osx" -f "$f" -o "$n.asm" -co "$n.c" 2> "$n.log"
done
# TODO: Do all the compilation here as well
