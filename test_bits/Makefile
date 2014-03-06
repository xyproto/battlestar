BITS:=$(shell getconf LONG_BIT)
GENERATED_ASM:=_hi.asm
GENERATED_C:=_hi.c

all: ../battlestarc hi$(BITS)

../battlestarc: ../main.go
	make -C ..

hi32:
	../battlestarc -bits=32 -f hi.bts -o $(GENERATED_ASM) -oc $(GENERATED_C) 2>/dev/null
	@hash gcc 2>/dev/null || { echo 'No gcc in path'; exit 1; }
	gcc -Os -m32 -nostdlib -c $(GENERATED_C) -o hi_c.o
	@hash yasm 2>/dev/null || { echo 'No yasm in path'; exit 1; }
	yasm -f elf32 -o hi_bts.o $(GENERATED_ASM)
	ld -s -melf_i386 --fatal-warnings -relax -o hi32 hi_bts.o hi_c.o
	hash sstrip 2>/dev/null && sstrip hi32 || strip -R .comment hi32
	@echo ok

hi64: hi32
	../battlestarc -bits=64 -f hi.bts -o $(GENERATED_ASM) -oc $(GENERATED_C) 2>/dev/null
	@hash gcc 2>/dev/null || { echo 'No gcc in path'; exit 1; }
	gcc -Os -m64 -nostdlib -c $(GENERATED_C) -o hi_c.o
	@hash yasm 2>/dev/null || { echo 'No yasm in path'; exit 1; }
	yasm -f elf64 -o hi_bts.o $(GENERATED_ASM)
	ld -s --fatal-warnings -relax -o hi64 hi_bts.o hi_c.o
	hash sstrip 2>/dev/null && sstrip hi64 || strip -R .comment hi64
	@echo ok

clean:
	-rm -f *.o $(GENERATED_ASM) $(GENERATED_C) hi32 hi64
