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
	@echo -en "\tGenerating $@ from $<... "
	@[ -e ../battlestar ] || echo 'Could not find ../battlestar!'
	@<$< ../battlestar -platform=$(LBITS) -osx=$(OSX) 2>$<.log >$@ && echo ok

%.o: %.asm
	@echo -en "\tCompiling $<... "
	#if [ $(OSX) = true ]; then
	@yasm -f elf$(LBITS) -o $@ $< && echo ok

%.elf: %.o
	@echo -en "\tLinking $@... "
	#if [ $(OSX) = true ]; then
	@ld -s --fatal-warnings -nostdlib --relax -o $@ $< && echo ok

clean:
	@echo 'Cleaning...'
	@rm -rf $(ASMSOURCES) $(OBJECTS) $(TARGETS) _obj *~ *.6 *.gch a.out *.log
