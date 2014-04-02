#!/bin/sh
echo yasm -fbin test04.asm -o test04
yasm -fbin test04.asm -o test04
dosbox -c "mount c ." -c "c:" -c cls -c test04.com -c exit > /dev/null
