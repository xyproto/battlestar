BITS:=$(shell getconf LONG_BIT)

all: clean
	<hi.bts ../battlestar --bits=$(BITS) 2>/dev/null > hi_bts.asm
	gcc -nostdlib -c hi.c -o hi_c.o
	yasm -f elf$(BITS) -o hi_bts.o hi_bts.asm
	ld -s --fatal-warnings -relax -o hi hi_bts.o hi_c.o

clean:
	-rm -f *.asm *.o hi
