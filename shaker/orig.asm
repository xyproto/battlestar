;shaker
;64 byte intro
;
;by frag
;http://fsqrt.blogspot.com/
;pzagrebin@gmail.com
;27-11-2011


;Code is pretty obvious -> no comments.

org 100h
	mov al, 13h
	int 10h
	push 0a000h
	push 09000h
	pop ds
	pop es

P:	mov [bx], ch	;Generates sine-like table
	add cx, si
	mov ax, 40
	imul cx
	sub si, dx
	inc bx
	jnz P

L:	mov si, 320
	mov ax, di
	xor dx, dx
	div si
	mov si, dx
	add al, byte [si+bx]
	mov si, ax
	add dl, byte [si+bx]
	xor al, dl
	stosb
	loop L
	
	mov al, [es:bx]
	mov dx, 3c9h
	out dx, al
	out dx, al
	out dx, al
	inc bx
	in al, 60h
	dec ax
	jnz L
	ret

	;Enjoy!


