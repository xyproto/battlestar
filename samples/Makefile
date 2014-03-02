LBITS := $(shell getconf LONG_BIT)
SOURCES := $(sort $(wildcard *.bts))
ASMSOURCES := $(SOURCES:.bts=.asm)
OBJECTS := $(ASMSOURCES:.asm=.o)

.SECONDARY: $(ASMSOURCES)

# OS X detection (depends on uname being present)
ifeq ($(shell uname -s),Darwin)
  OSX := true
  ASMCMD := yasm -f macho
  LDCMD := ld -macosx_version_min 10.8 -lSystem
  SUFFIX := macho
  LBITS = 32
else
  OSX := false
  ASMCMD := yasm -f elf$(LBITS)
  LDCMD := ld -s --fatal-warnings -nostdlib --relax
  SUFFIX := elf
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
	@[ -e ../battlestar ] || echo -e '$(ERROR_COLOR)Could not find ../battlestar!$(NO_COLOR)'
	@-< $< ../battlestar -bits=$(LBITS) -osx=$(OSX) 2> $<.log >$@

%.o: %.asm
	@[[ $< != *fail* ]] && (echo -e "  $(COMPILE_COLOR)$(BULLET) Compiling$(NO_COLOR) $<"; $(ASMCMD) -o $@ $<) || true

%.$(SUFFIX): %.o
	@[[ $< != *fail* ]] && (echo -e "  $(LINK_COLOR)$(BULLET) Linking$(NO_COLOR) $@\n"; $(LDCMD) -o $@ $<) || echo

clean:
	@echo Cleaning
	@rm -rf $(ASMSOURCES) $(OBJECTS) $(TARGETS) _obj *~ *.6 *.gch a.out *.log
