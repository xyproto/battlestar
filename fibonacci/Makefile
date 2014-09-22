all: clean main

main: main.c fib.o
	gcc -s -Os fib.o main.c -o main
	@-sstrip main 2>/dev/null || echo 'sstrip is not installed!'

fib.o: fib.bts
	bts compile fib.bts

clean:
	bts clean
	@rm -f main
