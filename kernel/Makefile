all:
	@echo
	bts build bootable boot.bts
	@mv boot boot.bin
	@chmod -x boot.bin
	@echo
	@echo 'Try running boot.bin with qemu (make boot)'
	@echo
	@echo 'NOTE: Work in progress, currently does not work as expected.'
	@echo '      Does not build on OS X, last time I tried (only Linux).'
	@echo

clean:
	@echo 'Cleaning...'
	bts clean
	@-rm -f *.asm *.c *.o *.log test?? boot linker.ld *.bin core.* a.out

boot:
	qemu-system-i386 -kernel boot.bin
