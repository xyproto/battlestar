#!/bin/sh

###########################################
#                                         #
#                                         #
# Run a .bts file as if it were a script  #
#                                         #
#                                         #
###########################################

# Check for the given source file
if [[ ! -f $1 ]]; then
  echo "No such file: $1"
  exit 1
fi

# Check for needed utilities
which battlestar >/dev/null || exit 1
which nasm >/dev/null || exit 1

# Set up temporary filenames
tmpfn1=`mktemp --suffix=asm`
tmpfn2=`mktemp --suffix=o`
tmpfn3=`mktemp --suffix=elf`

# Compile and link
BITS=`getconf LONG_BIT`
<$1 battlestar -platform=$BITS 2>/dev/null > "$tmpfn1"
nasm -f elf$BITS -o "$tmpfn2" "$tmpfn1"
ld -o "$tmpfn3" "$tmpfn2"

# Clean up after compiling and linking
rm "$tmpfn1" "$tmpfn2"

#echo "Size of executable: `du -b "$tmpfn3" | cut -d'/' -f1`bytes"

# Run the program
"$tmpfn3"

# Remove the program after execution
rm "$tmpfn3"
