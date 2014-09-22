.MODEL TINY
.386
.CODE
.STARTUP

        mov al, 13h
        int 10h

;        mov bh, 0A0h
;        mov es, bx
;        mov ds, bx

	pop sp 
	pop bx
	pop ds
	push ds
	pop es

        mov al, 62
        mov ch, 0FAh
        rep stosb


@@1:   
 	rol di, 3
        sub di, 7
        xor di, 2

        mov al, [di]
        add al, [di+321]
        shr al, 1

        push di

        stosb
        stosb
        add di, 013Eh
        stosb
        stosb

        pop di

	jmp @@1

;        in al, 60h
;        dec ax
;        jnz @@1

;        mov al, 03h
;        int 10h
;        ret

END