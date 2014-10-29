all: clean main

osx=$([[ `uname -s` = Darwin ]] && echo true || echo false)

main: clean main.c fib.o
	[ `uname -s` = Darwin ] && \
	  gcc -Os fib.o main.c -o main || \
	  gcc -s -Os fib.o main.c -o main
	@-sstrip main 2>/dev/null || echo 'sstrip is not installed!'

fib.o: fib.bts
	bts compile fib.bts -bits=64 -osx=true

clean:
	bts clean
	@rm -f main
