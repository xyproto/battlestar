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
	@echo
	@echo 'Booting ${BIN}...'
	@echo 'Press ctrl-alt-q to exit Qemu'
	@echo
	qemu-system-i386 -kernel ${BIN}

run: boot
