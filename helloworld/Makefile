all:
	@../scripts/generate.sh
	@# du -b does not work on OSX/BSD
	@echo "Only `ls -l hello | cut -d" " -f8` bytes!"

clean:
	@echo 'Cleaning...'
	@-rm -f *.asm *.c *.o *.log hello
