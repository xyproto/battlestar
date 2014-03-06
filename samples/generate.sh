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
cccmd="gcc -Os -m64 -nostdlib -nostdinc -std=c99"

if [ $bits = 32 ]; then
  asmcmd="yasm -f elf32"
  ldcmd='ld -s -melf_i386 --fatal-warnings -nostdlib --relax'
  cccmd='gcc -Os -m32 -nostdlib -nostdinc -std=c99'
fi

if [[ $1 == bootable ]]; then
  asmcmd="yasm -f elf32"
  ldcmd='ld -s -melf_i386 --fatal-warnings -nostdlib --relax'
  cccmd='gcc -Os -m32 -nostdlib -nostdinc -std=c99'
  echo 'Building a bootable kernel.'
  echo
  cccmd="$cccmd -ffreestanding -Wall -Wextra -fno-exceptions"
  echo "$cccmd"

  # From http://wiki.osdev.org/Bare_Bones#Linking_the_Kernel
  cat > linker.ld <<EOF
ENTRY(_start)

/* Tell where the various sections of the object files will be put in the final
   kernel image. */
SECTIONS
{
	/* Begin putting sections at 1 MiB, a conventional place for kernels to be
	   loaded at by the bootloader. */
	. = 1M;

	/* First put the multiboot header, as it is required to be put very early
	   early in the image or the bootloader won't recognize the file format.
	   Next we'll put the .text section. */
	.text BLOCK(4K) : ALIGN(4K)
	{
		*(.multiboot)
		*(.text)
	}

	/* Read-only data. */
	.rodata BLOCK(4K) : ALIGN(4K)
	{
		*(.rodata)
	}

	/* Read-write data (initialized) */
	.data BLOCK(4K) : ALIGN(4K)
	{
		*(.data)
	}

	/* Read-write data (uninitialized) and stack */
	.bss BLOCK(4K) : ALIGN(4K)
	{
		*(COMMON)
		*(.bss)
		*(.bootstrap_stack)
	}

	/* The compiler may produce other sections, by default it will put them in
	   a segment with the same name. Simply add stuff here as needed. */
}
EOF
  ldcmd="$ldcmd -T linker.ld"
  echo $ldcmd
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
  if [ -e $n ]; then
    require sstrip 2 && sstrip "$n" || strip -R .comment "$n"
  fi
  echo
done
