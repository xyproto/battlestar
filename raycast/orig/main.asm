; From https://www.pouet.net/prod.php?which=78044

mov al, 0x13
int 0x10

push 0x9FF6
pop es

mov dx,0x3c9

P:
  out dx, al
  out dx, al
  out dx, al
  cmp al, 63
  jz F
  inc ax
  F:
  loop P

pop ds

X:
  mov cl, -9
  L:
    mov   bl, cl
    mov   ax, 0xcccd
    mul   di
    lea   ax, [bx-80]
    add   al, dh
    imul  bl
    xchg  ax, dx   
    imul  bl
    mov   al, dh
    xor   al, ah
    sub   bl, [0x46c]
    add   al, 4
    and   al, bl
    test  al, 24
    loopz L
  or  al, 252
  sub al, cl
  stosb      
  loop X
