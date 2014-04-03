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

function abort {
  echo "$@"
  exit 1
}

function build {

  # Check if the .o files are intended to be used together with another compiler.
  # If so, clean up the generated source files and don't run the linker.
  other_compiler=false
  if [[ $1 = -c ]]; then
    other_compiler=true
    skipstrip=true
    shift
  fi

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
    battlestarc $params -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (cat "$n.log"; rm -f "$n.asm"; echo "$n.log" >> "$n.log"; echo "$n failed to build!"; return 1; )
  else
    battlestarc $params -f "$f" -o "$n.asm" -oc "$n.c" 2> "$n.log" || (rm -f "$n.asm"; echo "$n.log" >> "$n.log"; echo "$n failed to build (correct)"; return 2; )
  fi

  # Only return with an error code of the build failed and was not meant to fail
  retval=$?
  if [[ $retval == 2 ]]; then
    # Save the filenames for later cleaning
    echo -e "\n$n $n.asm $n.c $n.log" >> "$n.log"
    # Meant to fail, not a problem, return 0
    return 0
  elif [[ $retval != 0 ]]; then
    return $retval
  fi

  if [[ $bits = 16 ]]; then
    if [[ -f $n.c ]]; then
      echo "Skipping $f (inline C is not available for 16-bit x86)"
      # Save the filenames for later cleaning
      echo -e "\n$n.asm $n.c $n_c.o $n.o $n.log" >> "$n.log"
      # Quit
      if [[ $n == *fail* ]]; then
        # Meant to fail, ok
        return 0
      else
        # Not meant to fail, not ok
        return 1
      fi
    fi
  fi

  if [[ $linkfail = false ]]; then
    compiledc=false
    [ -e $n.c ] && ($cccmd -c "$n.c" -o "${n}_c.o" || abort "$n failed to compile")
    [ -e $n.c ] && compiledc=true
  else
    echo "WARNING: Can't compile inline C for 64-bit executables on a 32-bit system."
  fi
  [ -e $n.asm ] && ($asmcmd -o "$n.o" "$n.asm" || echo "$n failed to assemble")
  if [[ $other_compiler = true ]]; then
    # Clean up some of the files right away
    if [[ $compiledc = true ]]; then
      # Only remove the .c file if we are sure that we generated it
      rm -f "$n.c"
    fi
    # Remove the generated .asm file
    rm -f "$n.asm"
    # Save the filenames for later cleaning
    echo -e "\n$n.o ${n}_c.o $n.log" >> "$n.log"
    [ -e $n.o ] && return 0 || return 1
  fi
  if [[ $linkfail = false ]]; then
    if [[ $compiledc = true ]]; then
      $ldcmd "${n}_c.o" "$n.o" -o "$n" || echo "$n failed to link"
    elif [ -e $n.o ]; then
      if [[ $bits = 16 ]]; then
      	# The output file is a .com file
        mv "$n.o" "$n.com"
	# Create a script for running it with dosbox
	echo '#!/bin/sh' > "$n.sh"
	echo "dosbox -c \"mount c .\" -c \"c:\" -c cls -c $n.com -c pause -c exit > /dev/null" >> "$n.sh"
	chmod +x "$n.sh"
	# Save the filenames for later cleaning
	echo -e "\n$n.asm $n.c $n.com ${n}_c.o $n.log $n.sh" >> "$n.log"
	# Check if a .com file has been created and return a value accordingly
	[ -e $n.com ] && return 0 || return 1
      else
        $ldcmd "$n.o" -o "$n" || echo "$n failed to link"
      fi
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

  # Check if an executable has been generated and return a value accordingly
  [ -e $n ] && return 0 || return 1
}

# Should stripping be skipped?
skipstrip=false
#skipstrip=true

# Check for needed utilities
require battlestarc 1
require yasm 1
require ld 1
require gcc 0
require sstrip 0

# Discover if we are on a 32-bit or 64-bit system (bits is set to 32 or 64, or more?)
bits=`getconf LONG_BIT`

# Is it likely that gcc and ld will fail? (dealing with 64-bit executables on a 32-bit system)
linkfail=false

# Set bits if "bits=32", "bits=64" or "bits=16" is found in the arguments
if [[ $@ = *'bits=32'* ]]; then
  bits=32
elif [[ $@ = *'bits=64'* ]]; then
  if [[ $bits == 32 ]]; then
    linkfail=true
  fi
  bits=64
elif [[ $@ = *'bits=16'* ]]; then
  bits=16
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

if [ $bits = 16 ]; then
  asmcmd="yasm -f bin"
  ldcmd='ld -s --fatal-warnings -nostdlib --relax'
  cccmd="$stdgcc"
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
if [[ -f $2 ]]; then
  # For when -c is given
  build $@
elif [[ -f $1 ]]; then
  # For ordinary use
  build $@
else
  retval=0
  for f in *.bts; do
    if [[ -f $f ]]; then
      build $f $@
      retval=$?
      if [[ $retval != 0 ]]; then
	break
      fi
    fi
  done
  exit $retval
fi
