#!/bin/sh

function require {
  if [ $2 == 0 ]; then
    hash $1 2>/dev/null || { echo >&2 "Could not find $1 (optional)"; }
  elif [ $2 == 1 ]; then
    hash $1 2>/dev/null || { echo >&2 "Could not find $1. Aborting."; exit 1; }
  else
    hash $1 2>/dev/null || return 1
  fi
  return 0
}

# Check for needed utilities
require yasm 1
require ld 1
require gcc 0
require sstrip 0
echo

battlestarc=../battlestarc
if [ ! -e $battlestarc ]; then
  make -C ..
fi

bits=`getconf LONG_BIT`
osx=$([[ `uname -s` = Darwin ]] && echo true || echo false)
asmcmd="yasm -f elf64"
ldcmd='ld -s --fatal-warnings -nostdlib --relax'
stdgcc='gcc -Os -nostdlib -nostdinc -std=c99 -Wno-implicit -ffast-math -fno-inline -fomit-frame-pointer'
cccmd="$stdgcc -m64"

if [ $bits = 32 ]; then
  asmcmd="yasm -f elf32"
  ldcmd='ld -s -melf_i386 --fatal-warnings -nostdlib --relax'
  cccmd="$stdgcc -m32"
fi

skipstrip=false
if [[ $1 == bootable ]]; then
  echo 'Building a bootable kernel.'
  echo

  asmcmd="yasm -f elf32"
  echo $asmcmd

  cccmd="$stdgcc -m32 -ffreestanding -Wall -Wextra -fno-exceptions -Wno-implicit"
  echo "$cccmd"

  ldcmd='gcc -lgcc -nostdlib -Os -s -m32'
  if [ -e ../scripts/linker.ld ]; then
    ldcmd="$ldcmd -T ../scripts/linker.ld"
  elif [ -e linker.ld ]; then
    ldcmd="$ldcmd -T linker.ld"
  fi
  echo $ldcmd
  skipstrip=true
fi

if [[ $osx = true ]]; then
  asmcmd='yasm -f macho'
  ldcmd='ld -macosx_version_min 10.8 -lSystem'
  bits=32
fi

for f in *.bts; do
  n=`echo ${f/.bts} | sed 's/ //'`
  echo "Building $n"
  # Don't output the log if "fail" is in the filename
  if [[ $n != *fail* ]]; then
    $battlestarc -bits="$bits" -osx="$osx" -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (cat "$n.log"; rm -f "$n.asm"; echo "$n failed to build!")
  else
    $battlestarc -bits="$bits" -osx="$osx" -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (rm -f "$n.asm"; echo "$n failed to build (correct)")
  fi
  [ -e $n.c ] && ($cccmd -c "$n.c" -o "${n}_c.o" || echo "$n failed to compile")
  [ -e $n.asm ] && ($asmcmd -o "$n.o" "$n.asm" || echo "$n failed to assemble")
  if [ -e ${n}_c.o -a -e $n.o ]; then
    $ldcmd "${n}_c.o" "$n.o" -o "$n" || echo "$n failed to link"
  elif [ -e $n.o ]; then
    $ldcmd "$n.o" -o "$n" || echo "$n failed to link"
  fi
  if [[ $skipstrip == false ]]; then
    [ $osx = false ] && strip -R .comment -R .gnu.version "$n"
    require sstrip 2 && sstrip "$n"
  fi
  echo
done
