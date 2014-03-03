LBITS := $(shell getconf LONG_BIT)
SOURCES := $(sort $(wildcard *.bts))
ASMSOURCES := $(SOURCES:.bts=.asm)
OBJECTS := $(ASMSOURCES:.asm=.o)
BATTLESTAR := ../battlestar
$(shell ./generate.sh)

.SECONDARY: $(ASMSOURCES)

OSX := false
ASMCMD := yasm -f elf$(LBITS)
LDCMD := ld -s --fatal-warnings -nostdlib --relax
SUFFIX := elf
GCCCMD := gcc -Os -m64 -nostdlib

# 32-bit
ifeq ($(LBITS),32)
  LDCMD = ld -s -melf_i386 --fatal-warnings -nostdlib --relax
  GCCCMD = gcc -Os -m32 -nostdlib
endif

# OS X (depends on uname being present)
ifeq ($(shell uname -s),Darwin)
  OSX = true
  ASMCMD = yasm -f macho
  LDCMD = ld -macosx_version_min 10.8 -lSystem
  SUFFIX = macho
  LBITS = 32
endif

TARGETS = $(OBJECTS:.o=.$(SUFFIX))

all: clean $(TARGETS)

NO_COLOR=\x1b[0m
HEADER_COLOR=\x1b[31;01m
HEADER_COLOR2=\x1b[37;01m
GENERATE_COLOR=\x1b[35;01m
COMPILE_COLOR=\x1b[34;01m
LINK_COLOR=\x1b[32;01m
ERROR_COLOR=\x1b[31;01m
BULLET=->

%.asm: %.bts
	@echo -e "$(HEADER_COLOR)·---[$(HEADER_COLOR2) $< $(NO_COLOR)$(HEADER_COLOR)]-----------···$(NO_COLOR)\n"
	@echo -e "  $(GENERATE_COLOR)$(BULLET) Generating$(NO_COLOR) $@ from $<"
	@[ -e $(BATTLESTAR) ] || echo -e '$(ERROR_COLOR)Could not find $(BATTLESTAR)!$(NO_COLOR)'
	@-$(BATTLESTAR) -bits=$(LBITS) -osx=$(OSX) -f $< -o $@ -co $@.c 2> $<.log

%.o: %.asm
	@[[ $< != *fail* ]] && (echo -e "  $(COMPILE_COLOR)$(BULLET) Compiling assembly$(NO_COLOR) $<"; $(ASMCMD) -o $@ $<) || true
	@-[ -e $<.c ] && (echo -e "  $(COMPILE_COLOR)$(BULLET) Compiling C$(NO_COLOR) $<"; $(GCCCMD) -c $<.c -o $<.o) || true

%.$(SUFFIX): %.o
	# TODO: Find a general way to handle this, and also the files named *fail*
	@-[ -e test7.asm.o ] && (echo -e "  $(LINK_COLOR)$(BULLET) Linking$(NO_COLOR) $<"; $(LDCMD) -o $@ test7.asm.o test7.o)
	@-[ ! -e test7.asm.o ] && (echo -e "  $(LINK_COLOR)$(BULLET) Linking$(NO_COLOR) $@\n"; $(LDCMD) -o $@ $<) || echo

clean:
	@echo Cleaning
	@rm -rf $(ASMSOURCES) $(OBJECTS) $(TARGETS) _obj *~ *.6 *.gch a.out *.log *.asm.c
	$(shell ./clean.sh)
