LBITS := $(shell getconf LONG_BIT)
SOURCES := $(sort $(wildcard *.bts))
ASMSOURCES := $(SOURCES:.bts=.asm)
OBJECTS := $(ASMSOURCES:.asm=.o)
TARGETS := $(OBJECTS:.o=.elf)

all: $(TARGETS)

%.asm: %.bts
	@echo "--- $< ---"
	@echo -en "\tGenerating $@ from $<... "
	@<$< ../battlestar -platform=$(LBITS) 2>/dev/null > $@ && echo ok

%.o: %.asm
	@echo -en "\tCompiling $<... "
	@nasm -f elf$(LBITS) -o $@ $< && echo ok

%.elf: %.o
	@echo -en "\tLinking $@... "
	@ld --fatal-warnings -pie -nostdlib --relax -o $@ $< && echo ok

clean:
	@echo 'Cleaning...'
	@rm -rf _obj *~ *.6 *.gch a.out $(ASMSOURCES) $(OBJECTS) $(TARGETS)
