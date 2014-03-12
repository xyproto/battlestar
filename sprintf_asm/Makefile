CFLAGS=-m32 -g -O0
CC=gcc

all: test-sprinter test-csprinter test-itoa

test-sprinter: sprinter.s test-sprinter.c

test-itoa: itoa.s test-itoa.c

test-csprinter: csprinter.c test-sprinter.c
	$(CC) $(CFLAGS) -o test-csprinter csprinter.c test-sprinter.c

clean:
	rm -f test-sprinter test-csprinter test-itoa core
