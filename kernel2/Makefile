SRC=kernel.bts
BIN=kernel

all: ${BIN}

${BIN}: ${SRC}
	@echo
	bts build bootable ${SRC} --bits=32
	@chmod -x ${BIN}

clean:
	@echo 'Cleaning...'
	bts clean
	@-rm -f *.bin core.*

boot: ${BIN}
	@echo 'Booting...'
	qemu-system-x86_64 -kernel ${BIN}
