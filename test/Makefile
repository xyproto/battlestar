BITS:=$(shell getconf LONG_BIT)
GENERATED_ASM:=_hi.asm
GENERATED_C:=_hi.c

all: ../battlestar hi32 hi64

../battlestar: ../main.go
	make -C ..

hi32:
	../battlestar -bits=32 -f hi.bts -o $(GENERATED_ASM) -co $(GENERATED_C) 2>/dev/null
	gcc -Os -m32 -nostdlib -c $(GENERATED_C) -o hi_c.o
	yasm -f elf32 -o hi_bts.o $(GENERATED_ASM)
	ld -s -melf_i386 --fatal-warnings -relax -o hi32 hi_bts.o hi_c.o

hi64:
	../battlestar -bits=64 -f hi.bts -o $(GENERATED_ASM) -co $(GENERATED_C) 2>/dev/null
	gcc -Os -m64 -nostdlib -c $(GENERATED_C) -o hi_c.o
	yasm -f elf64 -o hi_bts.o $(GENERATED_ASM)
	ld -s --fatal-warnings -relax -o hi64 hi_bts.o hi_c.o

clean:
	-rm -f *.o $(GENERATED_ASM) $(GENERATED_C) hi32 hi64
