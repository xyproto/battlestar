.PHONY: todo clean

all: spongy.com

todo: spongy.com
	ndisasm original.com > original.asm
	ndisasm spongy.com > spongy.asm
	meld original.asm spongy.asm

spongy.bts: original.com
	com2bts original.com spongy.bts

spongy.com: spongy.bts
	bts build -bits=16
	bts size

clean:
	bts clean
	rm -f spongy.bts original.asm spongy.asm
