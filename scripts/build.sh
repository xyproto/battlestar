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

function build {
  f=$1
  shift
  params=$@
  echo "Building $f"
  n=`echo ${f/.bts} | sed 's/ //'`

  # TODO: This could probably be more robust
  if [[ $params != *bits* ]]; then
    params="$params -bits=$bits"
  fi

  # TODO: This could probably be more robust
  if [[ $params != *osx* ]]; then
    params="$params -osx=$osx"
  fi

  # Don't output the log if "fail" is in the filename
  if [[ $n != *fail* ]]; then
    #echo battlestarc $params -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (cat "$n.log"; rm -f "$n.asm"; echo "$n failed to build!")
    battlestarc $params -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (cat "$n.log"; rm -f "$n.asm"; echo "$n failed to build!")
  else
    #echo battlestarc $params -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (rm -f "$n.asm"; echo "$n failed to build (correct)")
    battlestarc $params -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (rm -f "$n.asm"; echo "$n failed to build (correct)")
  fi
  if [[ $linkfail = false ]]; then
    [ -e $n.c ] && ($cccmd -c "$n.c" -o "${n}_c.o" || echo "$n failed to compile")
  else
    echo "WARNING: Can't compile inline C for 64-bit executables on a 32-bit system."
  fi
  [ -e $n.asm ] && ($asmcmd -o "$n.o" "$n.asm" || echo "$n failed to assemble")
  if [[ $linkfail = false ]]; then
    if [ -e ${n}_c.o -a -e $n.o ]; then
      $ldcmd "${n}_c.o" "$n.o" -o "$n" || echo "$n failed to link"
    elif [ -e $n.o ]; then
      $ldcmd "$n.o" -o "$n" || echo "$n failed to link"
    fi
  else
    echo "WARNING: Can't link 64-bit executables on a 32-bit system."
  fi
  if [[ $skipstrip == false ]]; then
    [ $osx = false ] && (strip -R .comment -R .gnu.version "$n" 2>/dev/null)
    require sstrip 2 && (sstrip "$n" 2>/dev/null)
  fi
  # Save the filenames for later cleaning
  echo -e "\n$n $n.asm $n.c $n.o ${n}_c.o $n.log" >> "$n.log"
}

# Should stripping be skipped?
skipstrip=false
#skipstrip=true

# Check for needed utilities
require yasm 1
require ld 1
require gcc 0
require sstrip 0

# Discover if we are on a 32-bit or 64-bit system (bits is set to 32 or 64, or more?)
bits=`getconf LONG_BIT`

# Is it likely that gcc and ld will fail? (dealing with 64-bit executables on a 32-bit system)
linkfail=false

# Set bits if "bits=32" or "bits=64" is found in the arguments
if [[ $@ = *'bits=32'* ]]; then
  bits=32
elif [[ $@ = *'bits=64'* ]]; then
  if [[ $bits == 32 ]]; then
    linkfail=true
  fi
  bits=64
fi

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

if [[ $1 == bootable ]]; then
  shift

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

# Build one file or *.bts
if [[ -f $1 ]]; then
  build $@
else
  # TODO: If the first argument is a parameter, pass it on when building *.bts.
  #       If it looks like a filename, abort (already checked that it does not exist).
  #if [[ $1 != "" ]]; then
  #  echo 'Could not find $1!'
  #  exit 1
  #fi
  for f in *.bts; do
    if [[ -f $f ]]; then
      build $f $@
    fi
  done
fi
