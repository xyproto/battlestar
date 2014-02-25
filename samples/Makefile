LBITS := $(shell getconf LONG_BIT)

all: test1 test2

run:
	./run.sh test2.bts

test1.asm: test1.bts
	<test1.bts ../battlestar -platform=$(LBITS) 2>/dev/null > test1.asm

test2.asm: test2.bts
	<test2.bts ../battlestar -platform=$(LBITS) 2>/dev/null > test2.asm

test1.o: test1.asm
	nasm -f elf$(LBITS) -o test1.o test1.asm

test2.o: test2.asm
	nasm -f elf$(LBITS) -o test2.o test2.asm

test1: test1.o
	ld -o test1 test1.o

test2: test2.o
	ld -o test2 test2.o

clean:
	rm -rf _obj *.o *~ *.6 *.gch a.out test *.asm test1 test2
