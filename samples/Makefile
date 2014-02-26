LBITS := $(shell getconf LONG_BIT)
SOURCES := $(sort $(wildcard *.bts))
ASMSOURCES := $(SOURCES:.bts=.asm)
OBJECTS := $(ASMSOURCES:.asm=.o)
TARGETS := $(OBJECTS:.o=.elf)

.SECONDARY: $(ASMSOURCES)

# OS X detection (depends on uname being present)
OSX := false
ifeq ($(shell uname -s),Darwin)
  OSX = true
endif

# First clean, then build targets
all: clean $(TARGETS)

%.asm: %.bts
	@echo "--- $< ---"
	@echo -e "\tGenerating $@ from $<... "
	@[ -e ../battlestar ] || echo 'Could not find ../battlestar!'
	@<$< ../battlestar -platform=$(LBITS) -osx=$(OSX) 2>$<.log >$@

%.o: %.asm
	@echo -e "\tCompiling $<... "
	@[ $(OSX) = true ] && nasm -f macho -o $@ $< || yasm -f elf$(LBITS) -o $@ $<

%.elf: %.o
	@echo -e "\tLinking $@... "
	@[ $(OSX) = true ] && ld -s --fatal-warnings -nostdlib --relax -o $@ $< || ld -s --fatal-warnings -nostdlib --relax -o $@ $<

clean:
	@echo -e '\tCleaning...'
	@rm -rf $(ASMSOURCES) $(OBJECTS) $(TARGETS) _obj *~ *.6 *.gch a.out *.log
