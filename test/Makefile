BITS:=$(shell getconf LONG_BIT)

all: clean b64
	#b32

b32:
	<hi.bts ../battlestar --bits=32 2>/dev/null > hi_bts.asm
	gcc -nostdlib -c hi.c -o hi_c.o
	yasm -f elf32 -o hi_bts.o hi_bts.asm
	ld -s --fatal-warnings -relax -o hi hi_bts.o hi_c.o

b64:
	<hi.bts ../battlestar --bits=64 2>/dev/null > hi_bts.asm
	gcc -nostdlib -c hi.c -o hi_c.o
	yasm -f elf64 -o hi_bts.o hi_bts.asm
	ld -s --fatal-warnings -relax -o hi hi_bts.o hi_c.o

clean:
	-rm -f *.asm *.o hi
