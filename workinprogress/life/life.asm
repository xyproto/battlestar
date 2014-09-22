; Generated with Battlestar 0.3, at 2014-09-22 19:06

bits 16
org 0x100

section .text
jmp _start
;--- function main ---
_start:				; starting point of the program
main:				; name of the function


	mov al, 0x13		; al = 0x13
	;--- call interrupt 0x10 ---
	int 0x10			; perform the call

	pop sp			; asm

	pop bx				; stack -> bx

	pop ds			; asm

	push ds			; asm

	pop es			; asm

	mov al, 62		; al = 62
	mov ch, 0xFA		; ch = 0xFA
	rep stosb			; asm

	;--- loop ---
e_l1:					; start of loop e_l1

	rol di, 3			; asm

	sub di, 7			; di -= 7
	xor di, 2			; asm

	mov al, di			; asm

	add al, di+321			; asm

	shr al, 1		; al /= 2
	push di			; asm

	stosb			; asm label

	stosb			; asm label

	add di, 0x13E			; di += 0x13E
	stosb			; asm label

	stosb			; asm label

	pop di			; asm

	jmp e_l1				; loop forever
e_l1_end:				; end of loop e_l1
	;--- end of loop e_l1 ---


	;--- return from "main" ---
	ret			; exit program


