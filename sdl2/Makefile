CFLAGS=$(shell pkg-config sdl2 --cflags --libs)
LDFLAGS=$(shell pkg-config sdl2 --libs)

all:
	EXTERNLIB=1 bts build
	gcc -Os -nostdlib ${CFLAGS} main.o -o main
	strip main
	sstrip main
	bts size

clean:
	bts clean
