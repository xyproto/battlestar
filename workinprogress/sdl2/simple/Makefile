CC?=g++
#CC?=tcc     # Using tcc for linking results in even smaller executables
CFLAGS+=-pthread -fpie $(shell pkg-config sdl2 --cflags --libs)
LDFLAGS+=$(shell pkg-config sdl2 --libs)

all: clean
	@CFLAGS="${CFLAGS}" EXTERNLIB=1 bts build
	${CC} ${CFLAGS} -Os -nostdlib main.o -o main
	@-sstrip main 2>/dev/null || echo 'sstrip is not installed!'
	bts size

debug: clean
	@CFLAGS="-O1 -g" EXTERNLIB=1 bts build
	${CC} ${CFLAGS} -O1 -g -nostdlib main.o -o main
	bts size

clean:
	bts clean
