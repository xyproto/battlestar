BITS:=$(shell getconf LONG_BIT)
GENERATED_ASM:=_hi.asm
GENERATED_C:=_hi.c

all: ../battlestarc hi32 hi64

../battlestarc: ../main.go
	make -C ..

hi32:
	../battlestarc -bits=32 -f hi.bts -o $(GENERATED_ASM) -oc $(GENERATED_C) 2>/dev/null
	gcc -Os -m32 -nostdlib -c $(GENERATED_C) -o hi_c.o
	yasm -f elf32 -o hi_bts.o $(GENERATED_ASM)
	ld -s -melf_i386 --fatal-warnings -relax -o hi32 hi_bts.o hi_c.o
	-sstrip hi32

hi64:
	../battlestarc -bits=64 -f hi.bts -o $(GENERATED_ASM) -oc $(GENERATED_C) 2>/dev/null
	gcc -Os -m64 -nostdlib -c $(GENERATED_C) -o hi_c.o
	yasm -f elf64 -o hi_bts.o $(GENERATED_ASM)
	ld -s --fatal-warnings -relax -o hi64 hi_bts.o hi_c.o
	-sstrip hi64

clean:
	-rm -f *.o $(GENERATED_ASM) $(GENERATED_C) hi32 hi64
