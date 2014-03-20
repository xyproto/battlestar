all: boot

boot:
	@../scripts/generate.sh bootable
	@mv boot boot.bin
	@chmod -x boot.bin
	@echo 'Try running boot.bin with qemu'
	@echo

clean:
	@echo 'Cleaning...'
	@-rm -f *.asm *.c *.o *.log test?? boot linker.ld *.bin core.* a.out
