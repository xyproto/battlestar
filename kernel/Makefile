all:
	@echo
	bts build bootable boot.bts --bits=32
	@mv boot boot.bin
	@chmod -x boot.bin
	@cat README

clean:
	@echo 'Cleaning...'
	bts clean
	@-rm -f *.bin core.*

boot:
	qemu-system-i386 -kernel boot.bin
