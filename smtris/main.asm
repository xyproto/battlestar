; Generated with Battlestar 0.4, at 2016-02-08 12:14

bits 64
section .data
msg_author:	db "Written by Sebastian Mihai, 2014$" 		; constant string
_length_of_msg_author equ $ - msg_author	; size of constant value

msg_next:	db "Next$" 		; constant string
_length_of_msg_next equ $ - msg_next	; size of constant value

msg_left:	db "A - Left$" 		; constant string
_length_of_msg_left equ $ - msg_left	; size of constant value

msg_right:	db "S - Right$" 		; constant string
_length_of_msg_right equ $ - msg_right	; size of constant value

msg_rotate:	db "SPC - Rotate$" 		; constant string
_length_of_msg_rotate equ $ - msg_rotate	; size of constant value

msg_quit:	db "Q - Quit$" 		; constant string
_length_of_msg_quit equ $ - msg_quit	; size of constant value

msg_lines:	db "Lines$" 		; constant string
_length_of_msg_lines equ $ - msg_lines	; size of constant value

msg_game_over:	db "Game Over$" 		; constant string
_length_of_msg_game_over equ $ - msg_game_over	; size of constant value

msg_asmtris:	db "aSMtris$" 		; constant string
_length_of_msg_asmtris equ $ - msg_asmtris	; size of constant value

delay_centiseconds:	dq 5		; constant value
_length_of_delay_centiseconds equ $ - delay_centiseconds	; size of constant value

screen_width:	dq 320		; constant value
_length_of_screen_width equ $ - screen_width	; size of constant value

block_size:	dq 5		; constant value
_length_of_block_size equ $ - block_size	; size of constant value

block_per_piece:	dq 4		; constant value
_length_of_block_per_piece equ $ - block_per_piece	; size of constant value

colour_cemented_piece:	dq 0x40, 0x48, 0x54, 0x14, 0x42, 0x36, 0x34		; constant value
_length_of_colour_cemented_piece equ $ - colour_cemented_piece	; size of constant value

colour_falling_piece:	dq 39, 47, 55, 44, 6, 37, 33		; constant value
_length_of_colour_falling_piece equ $ - colour_falling_piece	; size of constant value

piece_t:	dq 1605, 1610, 1615, 3210, 10, 1610, 1615, 3210, 10, 1605, 1610, 1615, 10, 1605, 1610, 3210		; constant value
_length_of_piece_t equ $ - piece_t	; size of constant value

piece_j:	dq 1605, 1610, 1615, 3215, 10, 15, 1610, 3210, 5, 1605, 1610, 1615, 10, 1610, 3205, 3210		; constant value
_length_of_piece_j equ $ - piece_j	; size of constant value

piece_l:	dq 1605, 1610, 1615, 3205, 10, 1610, 3210, 3215, 15, 1605, 1610, 1615, 5, 10, 1610, 3210		; constant value
_length_of_piece_l equ $ - piece_l	; size of constant value

piece_z:	dq 1605, 1610, 3210, 3215, 15, 1610, 1615, 3210, 1605, 1610, 3210, 3215, 15, 1610, 1615, 3210		; constant value
_length_of_piece_z equ $ - piece_z	; size of constant value

piece_s:	dq 1610, 1615, 3205, 3210, 10, 1610, 1615, 3215, 1610, 1615, 3205, 3210, 10, 1610, 1615, 3215		; constant value
_length_of_piece_s equ $ - piece_s	; size of constant value

piece_square:	dq 1605, 1610, 3205, 3210, 1605, 1610, 3205, 3210, 1605, 1610, 3205, 3210, 1605, 1610, 3205, 3210		; constant value
_length_of_piece_square equ $ - piece_square	; size of constant value

piece_line:	dq 1600, 1605, 1610, 1615, 10, 1610, 3210, 4810, 1600, 1605, 1610, 1615, 10, 1610, 3210, 4810		; constant value
_length_of_piece_line equ $ - piece_line	; size of constant value

section .text
global _start			; make label available to the linker
_start:				; starting point of the program



	pieces_origin:			; asm label

	;--- exit program ---
	mov rax, 60			; function call: 60
	xor rdi, rdi			; return code 0
	syscall				; exit program


