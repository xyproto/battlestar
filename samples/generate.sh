#!/bin/sh

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
  $battlestar -bits="$bits" -osx="$osx" -f "$f" -o "$n.asm" -co "$n.c" 2> "$n.log"
  [ -e $n.c ] && $gcccmd -c "$n.c" -o "${n}_c.o"
  [ -e $n.asm ] && $asmcmd -o "$n.o" "$n.asm"
  if [ -e ${n}_c.o -a -e $n.o ]; then
    $ldcmd "${n}_c.o" "$n.o" -o "$n" || echo "$n failed"
  elif [ -e $n.o ]; then
    $ldcmd "$n.o" -o "$n" || echo "$n failed"
  fi
done
