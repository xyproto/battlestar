CC?=gcc
#CC?=tcc     # Using tcc for linking results in even smaller executables
CFLAGS+=$(shell pkg-config sdl2 --cflags --libs)
LDFLAGS+=$(shell pkg-config sdl2 --libs)

all: clean
	@CFLAGS="${CFLAGS}" EXTERNLIB=1 bts build
	${CC} -Os -nostdlib main_c.o main.o ${CFLAGS} -o main
	-sstrip main
	bts size

debug: clean
	@CFLAGS="-O1 -g" EXTERNLIB=1 bts build
	gcc -O1 -g -nostdlib main_c.o main.o ${CFLAGS} -o main
	bts size

clean:
	bts clean
