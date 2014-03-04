#!/bin/sh

# Check for needed utilities
which yasm >/dev/null || (echo Could not find yasm; exit 1)

battlestar=../battlestar
if [ ! -e $battlestar ]; then
  make -C ..
fi

bits=`getconf LONG_BIT`
osx=$([[ `uname -s` = Darwin ]] && echo true || echo false)
asmcmd="yasm -f elf$bits"
ldcmd='ld -s --fatal-warnings -nostdlib --relax'
gcccmd="gcc -Os -m64 -nostdlib"

if [[ $bits = 32 ]]; then
  ldcmd='ld -s -melf_i386 --fatal-warnings -nostdlib --relax'
  gcccmd='gcc -Os -m32 -nostdlib'
fi

if [[ $osx = true ]]; then
  asmcmd='yasm -f macho'
  ldcmd='ld -macosx_version_min 10.8 -lSystem'
  bits=32
fi

for f in *.bts; do
  n=`echo ${f/.bts} | sed 's/ //'`
  echo "Building $n"
  # TODO: Don't output the log if "fail" is in the filename
  $battlestar -bits="$bits" -osx="$osx" -f "$f" -o "$n.asm" -co "$n.c" 2> "$n.log" || (cat "$n.log"; rm "$n.asm"; echo "$n failed to build.")
  [ -e $n.c ] && ($gcccmd -c "$n.c" -o "${n}_c.o" || echo "$n failed to compile")
  [ -e $n.asm ] && ($asmcmd -o "$n.o" "$n.asm" || echo "$n failed to assemble")
  if [ -e ${n}_c.o -a -e $n.o ]; then
    $ldcmd "${n}_c.o" "$n.o" -o "$n" || echo "$n failed to link"
  elif [ -e $n.o ]; then
    $ldcmd "$n.o" -o "$n" || echo "$n failed to link"
  fi
done
