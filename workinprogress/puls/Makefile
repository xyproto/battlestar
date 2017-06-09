.PHONY: todo clean

all: puls.com

compare: puls.com
	ndisasm -b 16 original.com > original.asm
	ndisasm -b 16 puls.com > puls.asm
	meld original.asm puls.asm

puls.bts: original.com
	com2bts original.com puls.bts

puls.com: puls.bts
	bts build -bits=16
	bts size

clean:
	bts clean
	rm -f puls.bts original.asm puls.asm
