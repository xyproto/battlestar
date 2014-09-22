all: clean main

main: main.c fib.o
	gcc -s -Os fib.o main.c -o main
	-sstrip main

fib.o: fib.bts
	bts compile fib.bts

clean:
	bts clean
	@rm -f main
