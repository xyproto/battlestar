NAME=ok
BITS:=$(shell getconf LONG_BIT)
GENERATED_ASM:=_$(NAME).asm
GENERATED_C:=_$(NAME).c

all: ../battlestarc $(NAME)$(BITS)

../battlestarc: ../main.go
	make -C ..

$(NAME)32:
	../battlestarc -bits=32 -f $(NAME).bts -o $(GENERATED_ASM) -oc $(GENERATED_C) 2>/dev/null
	@hash yasm 2>/dev/null || { echo 'No yasm in path'; exit 1; }
	yasm -f elf32 -o $(NAME)_bts.o $(GENERATED_ASM)
	ld -s -melf_i386 --fatal-warnings -relax -o $(NAME)32 $(NAME)_bts.o
	hash sstrip 2>/dev/null && sstrip $(NAME)32 || strip -R .comment $(NAME)32
	@echo ok

$(NAME)64: $(NAME)32
	../battlestarc -bits=64 -f $(NAME).bts -o $(GENERATED_ASM) -oc $(GENERATED_C) 2>/dev/null
	@hash yasm 2>/dev/null || { echo 'No yasm in path'; exit 1; }
	yasm -f elf64 -o $(NAME)_bts.o $(GENERATED_ASM)
	ld -s --fatal-warnings -relax -o $(NAME)64 $(NAME)_bts.o
	hash sstrip 2>/dev/null && sstrip $(NAME)64 || strip -R .comment $(NAME)64
	@echo ok

clean:
	-rm -f *.o $(GENERATED_ASM) $(GENERATED_C) $(NAME)32 $(NAME)64
