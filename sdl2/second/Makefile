CFLAGS+=$(shell pkg-config sdl2 --cflags --libs)
LDFLAGS+=$(shell pkg-config sdl2 --libs)

all: clean
	EXTERNLIB=1 bts build
	@#tcc -Os -nostdlib main.o ${CFLAGS} -o main || gcc -Os -nostdlib main.o ${CFLAGS} -o main
	gcc -Os -nostdlib main.o ${CFLAGS} -o main
	strip main
	sstrip main
	bts size

debug: clean
	CFLAGS="-O1 -g" EXTERNLIB=1 bts build
	gcc -O0 -g -nostdlib main.o ${CFLAGS} -o main
	bts size

clean:
	bts clean
