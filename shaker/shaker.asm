; Generated with Battlestar 0.4, at 2014-11-07 18:25

bits 16
org 0x100

section .text
_start:				; starting point of the program

	mov al, 0x13		; al = 0x13
	;--- call interrupt 0x10 ---
	int 0x10			; perform the call

	push 0x0a000			; asm

	push 0x09000			; asm

	pop ds				; stack -> ds

	pop es				; stack -> es

	P:			; asm label

	mov bx, ch			; asm

	add cx, si			; cx += si
	mov ax, 40		; ax = 40
	imul cx			; asm

	sub si, dx			; si -= dx
	inc bx			; bx++
	jnz P			; asm

	L:			; asm label

	mov si, 320		; si = 320
	mov ax, di			; ax = di
	xor dx, dx		; dx = 0
	div si			; asm

	mov si, dx			; si = dx
	add al, [si+bx]			; al += [si+bx]
	mov si, ax			; si = ax
	add dl, [si+bx]			; dl += [si+bx]
	xor al, dl			; al ^= dl
	stosb			; write the value in al, starting at es:di

	loop L			; asm

	mov al, [es:bx]			; asm

	mov dx, 0x3c9		; dx = 0x3c9
	out dx, al			; asm

	out dx, al			; asm

	out dx, al			; asm

	inc bx			; bx++
	in al, 60h			; asm

	dec ax			; ax--
	jnz L			; asm

	;--- exit program ---
	mov ah, 0x4c			; function 4C
	xor al, al			; exit code 0
	int 0x21			; exit program


