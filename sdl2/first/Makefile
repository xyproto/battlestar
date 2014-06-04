CFLAGS=$(shell pkg-config sdl2 --cflags --libs)
LDFLAGS=$(shell pkg-config sdl2 --libs)

all:
	EXTERNLIB=1 bts build
	gcc -Os -nostdlib main.o ${CFLAGS} -o main
	@#tcc -Os -nostdlib main.o ${CFLAGS} -o main || gcc -Os -nostdlib main.o ${CFLAGS} -o main
	strip main
	sstrip main
	bts size

clean:
	bts clean
