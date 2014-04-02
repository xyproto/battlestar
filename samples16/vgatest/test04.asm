; Generated with Battlestar 0.1, at 2014-04-02 14:50

bits 16
section .text
org 0x100
jmp start

;--- function wait_for_keypress ---
wait_for_keypress:				; name of the function


	mov ah, 0x10		; ah = 0x10
	;--- call interrupt 0x16 ---
	int 0x16			; perform the call


	;--- return from "wait_for_keypress" ---
	ret				; Return

;--- function graphics_mode ---
graphics_mode:				; name of the function


	xor ah, ah		; ah = 0
	mov al, 0x13		; al = 0x13
	;--- call interrupt 0x10 ---
	int 0x10			; perform the call


	;--- return from "graphics_mode" ---
	ret				; Return

;--- function text_mode ---
text_mode:				; name of the function


	xor ah, ah		; ah = 0
	mov al, 0x03		; al = 0x03
	;--- call interrupt 0x10 ---
	int 0x10			; perform the call


	;--- return from "text_mode" ---
	ret				; Return

;--- function main ---
start:				; starting point of the program

	;--- call the "graphics_mode" function ---
	call graphics_mode

	push 0xa000
	pop es

	xor di, di
	mov ax, 0x2727
	mov cx, 0x7d00

	rep stosw

	;--- call the "wait_for_keypress" function ---
	call wait_for_keypress

	;--- call the "text_mode" function ---
	call text_mode


	;--- return from "main" ---
	ret			; exit program



