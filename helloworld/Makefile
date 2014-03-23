all:
	@../scripts/generate.sh
	@# du -b does not work on OSX/BSD
	@[ `uname -s` = Linux ] && echo "Only `du -b hello | cut -f1 ` bytes!" || echo "Only `ls -l hello | cut -d" " -f8` bytes!"

clean:
	@echo 'Cleaning...'
	@-rm -f *.asm *.c *.o *.log hello
