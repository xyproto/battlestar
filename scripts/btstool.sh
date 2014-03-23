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

function require {
  if [ $2 == 0 ]; then
    hash $1 2>/dev/null || { echo >&2 "Could not find $1 in path (optional)."; }
  elif [ $2 == 1 ]; then
    hash $1 2>/dev/null || { echo >&2 "Could not find $1 in path. Aborting."; exit 1; }
  else
    hash $1 2>/dev/null || return 1
  fi
  return 0
}

# Check for needed utilities
require "$battlestarc" 1
require yasm 1
require ld 1
require gcc 2
require sstrip 2

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

    # Strip the program, if available
    require sstrip 2 && sstrip "$elffn"
  
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
  shift
  btsbuild $@
  exit 0
fi

if [[ $1 == clean ]]; then
  # For each log file
  for log in *.log; do
    if [[ -f $log ]]; then
      # Remove all the files mentioned at the end
      for fn in `tail -1 $log`; do
	if [[ -f $fn ]]; then
	  rm -fv "$fn"
	fi
      done
    fi
  done
  # Clean *~ and a.out as well
  rm -fv a.out *~
  exit 0
fi

if [[ $1 == size ]]; then
  # For each log file, print the size of the resulting executable
  for log in *.log; do
    n=`echo ${log/.log} | sed 's/ //'`
    if [[ -f $n ]]; then
      # du -b does not work on OSX/BSD
      [ `uname -s` = Linux ] && echo "$n is `du -b $n | cut -f1 ` bytes" || echo "$n is `ls -l $n | cut -d" " -f8` bytes"
    fi
  done
  exit 0
fi

# Unknown command, assume run, for ease of use of #!/usr/bin/bts at top of scripts
[[ -e $1 ]] && run $1 || run $2
