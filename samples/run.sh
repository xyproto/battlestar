#!/bin/sh
if [[ ! -f $1 ]]; then
  echo "No such file: $1"
  exit 1
fi
<$1 ../battlestar 2>/dev/null > /tmp/tmp.asm
nasm -f elf64 -o /tmp/tmp.o /tmp/tmp.asm
rm /tmp/tmp.asm
ld -o /tmp/tmp.elf /tmp/tmp.o
rm /tmp/tmp.o
/tmp/tmp.elf
rm /tmp/tmp.elf
