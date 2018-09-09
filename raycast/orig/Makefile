.PHONY: all clean

all: main.com

main.com: main.asm
	yasm -f bin -o main.com main.asm

run: main.com
	dosbox main.com

clean:
	rm -f main.com
