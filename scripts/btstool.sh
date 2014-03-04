#!/bin/sh

###########################################
#                                         #
#                                         #
# Run a .bts file as if it were a script  #
#                                         #
#                                         #
###########################################

# Check for needed utilities (in PATH)
battlestarc=battlestarc
which "$battlestarc" 2>/dev/null || echo 'Could not find battlestar compiler'
which "$battlestarc" 2>/dev/null || exit 1
which yasm 2>/dev/null || echo 'Could not find yasm'
which yasm 2>/dev/null || exit 1
which gcc 2>/dev/null || echo 'Could not find gcc (optional)'

bits=`getconf LONG_BIT`
osx=$([[ `uname -s` = Darwin ]] && echo true || echo false)
asmcmd="yasm -f elf$bits"
ldcmd='ld -s --fatal-warnings -nostdlib --relax'
cccmd="gcc -Os -m64 -nostdlib"

if [[ $bits = 32 ]]; then
  ldcmd='ld -s -melf_i386 --fatal-warnings -nostdlib --relax'
  cccmd='gcc -Os -m32 -nostdlib'
fi

if [[ $osx = true ]]; then
  asmcmd='yasm -f macho'
  ldcmd='ld -macosx_version_min 10.8 -lSystem'
  bits=32
fi

function run {
  # Check for the given source file
  if [[ ! -f $1 ]]; then
    echo "No such file: $1"
    exit 1
  fi
  
  # Set up temporary filenames
  asmfn=`mktemp --suffix=.asm`
  o1fn=`mktemp --suffix=.o`
  o2fn=`mktemp --suffix=.o`
  elffn=`mktemp --suffix=.elf`
  cfn=`mktemp --suffix=.c`
  logfn=`mktemp --suffix=.log`
  
  # Compile and link
  $battlestarc -bits="$bits" -osx="$osx" -f $1 -o "$asmfn" -oc "$cfn" 2>"$logfn" || (cat "$logfn"; rm "$asmfn"; echo "$1 failed to build.")
  if [ -e "$asmfn" ]; then
    [ -e $cfn ] && ($cccmd -c "$cfn" -o "${o2fn}" || echo "$1 failed to compile")
    [ -e $asmfn ] && ($asmcmd -o "$o1fn" "$asmfn" || echo "$1 failed to assemble")
    if [ -e ${o2fn} -a -e ${o1fn} ]; then
      $ldcmd "${o2fn}" "$o1fn" -o "$elffn" || echo "$1 failed to link"
    elif [ -e ${o1fn} ]; then
      $ldcmd "${o1fn}" -o "$elffn" || echo "$1 failed to link"
    fi
  
    # Clean up after compiling and linking
    rm -f "${o1fn}" "${o2fn}" "$asmfn" "$cfn" "$logfn"
  
    #echo
    #echo "Running $1"
    #echo "Size of executable: `du -b "$elffn" | cut -d'/' -f1`bytes"
    #echo
  
    # Run the program
    [ -e $elffn ] && "$elffn"
  
    # Remove the program after execution
    [ -e $elffn ] && rm "$elffn"
  fi
}

# Check for a command
if [[ $1 == run ]]; then
  run $2
  exit 0
fi

if [[ $1 == build ]]; then
  echo To implement: build
  exit 1
fi

# Unknown command, assume run, for ease of use of #!/usr/bin/bts at top of scripts
[[ -e $1 ]] && run $1 || run $2
