CFLAGS=$(shell pkg-config sdl2 --cflags --libs)
LDFLAGS=$(shell pkg-config sdl2 --libs)

all: clean
	EXTERNLIB=1 bts build
	tcc -Os -nostdlib main.o ${CFLAGS} -o main || gcc -Os -nostdlib main.o ${CFLAGS} -o main
	strip main
	sstrip main
	bts size

clean:
	bts clean
