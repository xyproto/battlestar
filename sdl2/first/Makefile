CFLAGS=$(shell pkg-config sdl2 --cflags --libs)
LDFLAGS=$(shell pkg-config sdl2 --libs)

all:
	EXTERNLIB=1 bts build
	echo hi
	tcc -Os -nostdlib main.o ${CFLAGS} -o main || gcc -Os -nostdlib main.o ${CFLAGS} -o main
	@#gcc -Os main.o -o main ${CFLAGS}
	@#ld main.o -o main ${LDFLAGS}
	strip main
	sstrip main
	bts size

clean:
	bts clean
