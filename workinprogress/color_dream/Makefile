.PHONY: todo clean

all: color_dream.com

compare: color_dream.com
	ndisasm original.com > original.asm
	ndisasm color_dream.com > color_dream.asm
	meld original.asm color_dream.asm

color_dream.bts: original.com
	com2bts original.com color_dream.bts

color_dream.com: color_dream.bts
	bts build -bits=16
	bts size

clean:
	bts clean
	rm -f color_dream.bts original.asm color_dream.asm
