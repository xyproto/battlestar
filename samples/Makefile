LBITS := $(shell getconf LONG_BIT)
SOURCES := $(sort $(wildcard *.bts))
ASMSOURCES := $(SOURCES:.bts=.asm)
OBJECTS := $(ASMSOURCES:.asm=.o)
TARGETS := $(OBJECTS:.o=.elf)

.SECONDARY: $(ASMSOURCES)

all: $(TARGETS)

%.asm: %.bts
	@echo "--- $< ---"
	@echo -en "\tGenerating $@ from $<... "
	@[ -e ../battlestar ] || echo Could not find ../battlestar! && exit 1
	@<$< ../battlestar -platform=$(LBITS) 2> $<.log > $@ && echo ok

%.o: %.asm
	@echo -en "\tCompiling $<... "
	@nasm -f elf$(LBITS) -o $@ $< && echo ok

%.elf: %.o
	@echo -en "\tLinking $@... "
	@ld -s --fatal-warnings -nostdlib --relax -o $@ $< && echo ok

clean:
	@echo 'Cleaning...'
	@rm -rf $(ASMSOURCES) $(OBJECTS) $(TARGETS) _obj *~ *.6 *.gch a.out *.log
