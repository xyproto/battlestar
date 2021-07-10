#!/bin/bash

###########################################
#                                         #
#                                         #
# Run a .bts file as if it were a script  #
#                                         #
#                                         #
###########################################

function require {
  if [ $2 == 0 ]; then
    hash $1 2>/dev/null || ( echo >&2 "Could not find $1 in path (optional)."; )
  elif [ $2 == 1 ]; then
    hash $1 2>/dev/null || ( echo >&2 "Could not find $1 in path. Aborting."; exit 1; )
  else
    hash $1 2>/dev/null || return 1
  fi
  return 0
}

require uname 1

# Select yasm or nasm
hash yasm 2>/dev/null && assembler=yasm || assembler=nasm

bits=`getconf LONG_BIT`
osx=$([[ `uname -s` = Darwin ]] && echo true || echo false)
asmcmd="$assembler -f elf$bits"
ldcmd="ld -s --fatal-warnings -nostdlib --relax"
cccmd="gcc -Os -std=c99 -m64 -nostdlib $CFLAGS"

if [[ $bits = 32 ]]; then
  ldcmd="ld -s -melf_i386 --fatal-warnings -nostdlib --relax"
  cccmd="gcc -Os -std=c99 -nostdlib -Wno-impplicit -ffast-math -fno-inline -fomit-frame-pointer -m32 $CFLAGS"
fi

if [[ $osx = true ]]; then
  asmcmd='$assembler -f macho'
  ldcmd="ld -macosx_version_min 10.8 -lSystem $LDFLAGS"
  #bits=32
fi

function asm64check {
  if [[ $bits = 64 ]]; then
    if [[ $assembler = nasm ]]; then
      echo "Nasm does not support 64-bit assembly correctly. Please install Yasm."
      exit 1
    fi
  fi
}

function usage {
  echo 'Battlestar'
  echo
  echo 'Arguments:'
  echo
  echo ' bts [FILE]                    - run the given file as a script'
  echo ' bts run [FILE]                - run the given file as a script'
  echo ' bts build [FILE]              - build native executable'
  echo ' bts build --bits=64 [FILE]    - build native 64-bit executable'
  echo ' bts build --bits=32 [FILE]    - build native 32-bit executable'
  echo ' bts build --bits=16 [FILE]    - build native 16-bit executable'
  echo '                                 also create dosbox launcher script'
  echo ' bts compile [FILE]            - build object file'
  echo ' bts clean                     - remove stray files'
  echo ' bts size                      - analyze log files after building'
  echo ' bts help                      - this help'
  echo
}

function run {
  # Check for the given source file
  if [[ -z $1 ]]; then
    usage
    exit 1
  fi
  if [[ $1 == "--help" ]]; then
    usage
    exit 1
  fi
  if [[ $1 == "help" ]]; then
    usage
    exit 1
  fi
  if [[ ! -f $1 ]]; then
    echo "No such file: $1"
    exit 1
  fi

  asm64check

  # Set up temporary filenames
  asmfn=`mktemp --suffix=.asm`
  o1fn=`mktemp --suffix=.o`
  o2fn=`mktemp --suffix=.o`
  elffn=`mktemp --suffix=.elf`
  cfn=`mktemp --suffix=.c`
  logfn=`mktemp --suffix=.log`

  # Compile and link
  if ! battlestarc -bits="$bits" -osx="$osx" -f $1 -o "$asmfn" -oc "$cfn" 2>"$logfn"; then
    cat "$logfn"
    rm "$asmfn"
    echo "$1 failed to build."
    exit 1
  fi
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
    retval=$?

    # Remove the program after execution
    [ -e $elffn ] && rm "$elffn"

    # Exit with the same exit code as the program
    exit $retval
  fi
}

# The "run" command
if [[ $1 == run ]]; then
  # Check for needed utilities
  require battlestarc 1
  require $assembler 1
  require ld 1
  require gcc 2
  require sstrip 2
  asm64check

  shift
  run $@
  exit $?
fi

# The "build" command
if [[ $1 == build ]]; then
  require btsbuild 1

  if [[ ! $2 == "-bits=16" ]] && [[ ! $2 == "-bits=32" ]]; then
    asm64check
  fi

  shift
  btsbuild $@
  exit $?
fi

# The "compile" command
if [[ $1 == compile ]]; then
  require btsbuild 1
  asm64check

  shift
  btsbuild -c $@
  exit $?
fi

# The "clean" command
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

# The "size" command
if [[ $1 == size ]]; then
  # For each log file, print the size of the resulting executable
  for log in *.log; do
    n=`echo ${log/.log} | sed 's/ //'`
    if [[ -f $n.com ]]; then
      # du -b does not work on OSX/BSD
      [ `uname -s` = Linux ] && echo "$n.com is `du -b $n.com | cut -f1 ` bytes" || echo "$n.com is `ls -l $n.com | cut -d" " -f8` bytes"
    elif [[ -f $n ]]; then
      # du -b does not work on OSX/BSD
      [ `uname -s` = Linux ] && echo "$n is `du -b $n | cut -f1 ` bytes" || echo "$n is `ls -l $n | cut -d" " -f8` bytes"
    fi
  done
  exit $?
fi

# Unknown command, assume run, for ease of use of #!/usr/bin/bts at top of scripts
[[ -e $1 ]] && run $1 || run $2
