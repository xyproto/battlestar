; Generated with Battlestar 0.3, at 2014-04-09 17:33

bits 64
section .data
bottle:	db " bottle" 		; constant string
_length_of_bottle equ $ - bottle	; size of constant value

plural:	db "s" 		; constant string
_length_of_plural equ $ - plural	; size of constant value

ofbeer:	db " of beer" 		; constant string
_length_of_ofbeer equ $ - ofbeer	; size of constant value

wall:	db " on the wall" 		; constant string
_length_of_wall equ $ - wall	; size of constant value

sep:	db ", " 		; constant string
_length_of_sep equ $ - sep	; size of constant value

takedown:	db "Take one down and pass it around, " 		; constant string
_length_of_takedown equ $ - takedown	; size of constant value

u_no:	db "No" 		; constant string
_length_of_u_no equ $ - u_no	; size of constant value

l_no:	db "no" 		; constant string
_length_of_l_no equ $ - l_no	; size of constant value

more:	db " more bottles of beer" 		; constant string
_length_of_more equ $ - more	; size of constant value

store:	db "Go to the store and buy some more, " 		; constant string
_length_of_store equ $ - store	; size of constant value

dotnl:	db ".", 10 		; constant string
_length_of_dotnl equ $ - dotnl	; size of constant value

nl:	db "", 10 		; constant string
_length_of_nl equ $ - nl	; size of constant value

section .text
;--- function writenum ---
global writenum			; make label available to the linker
writenum:				; name of the function

	;--- setup stack frame ---
	push rbp			; save old base pointer
	mov rbp, rsp			; use stack pointer as new base pointer

	mov rbx, rax			; rbx = rax
	;--- loop ---
e_l1:					; start of loop e_l1

	cmp rax, 10			; compare
	jl e_l1_end			; break


	;--- signed division: rax /= 10 ---
	xor rdx, rdx		; rdx = 0 (64-bit 0:rax instead of 128-bit rdx:rax)
	mov r8, 10		; divisor, r8 = 10
	idiv r8			; rax = rdx:rax / r8

	mov rbx, rdx			; rbx = rdx
	add rax, 48			; rax += 48
	;--- system call ---
	sub rsp, 8			; make some space for storing rax on the stack
	mov QWORD [rsp], rax		; move rax to a memory location on the stack
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, rsp			; parameter #2 is rsp
	mov rdx, 4			; parameter #3 is 4
	syscall				; perform the call
	add rsp, 8			; move the stack pointer back

	jmp e_l1_end			; break

	jmp e_l1				; loop forever
e_l1_end:				; end of loop e_l1
	;--- end of loop e_l1 ---

	mov rax, rbx			; rax = rbx
	add rax, 48			; rax += 48
	;--- system call ---
	sub rsp, 8			; make some space for storing rax on the stack
	mov QWORD [rsp], rax		; move rax to a memory location on the stack
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, rsp			; parameter #2 is rsp
	mov rdx, 4			; parameter #3 is 4
	syscall				; perform the call
	add rsp, 8			; move the stack pointer back

	;--- takedown stack frame ---
	mov rsp, rbp			; use base pointer as new stack pointer
	pop rbp				; get the old base pointer


	;--- return from "writenum" ---
	ret				; Return

;--- function pluralifneeded ---
global pluralifneeded			; make label available to the linker
pluralifneeded:				; name of the function

	;--- setup stack frame ---
	push rbp			; save old base pointer
	mov rbp, rsp			; use stack pointer as new base pointer

	;--- if1 ---
	cmp rax, 1			; compare
	jle if1_end			; break

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, plural		; parameter #2 is &plural
	mov rdx, _length_of_plural		; parameter #3 is len(plural)
	syscall				; perform the call

if1_end:				; end of if block if1

	;--- takedown stack frame ---
	mov rsp, rbp			; use base pointer as new stack pointer
	pop rbp				; get the old base pointer


	;--- return from "pluralifneeded" ---
	ret				; Return

;--- function main ---
global main			; make label available to the linker
global _start			; make label available to the linker
_start:				; starting point of the program
main:				; name of the function


	;--- loop 99 times ---
	mov rcx, 99			; initialize loop counter
l2:					; start of loop l2
	push rcx			; save the counter

	push cx			; cx -> stack

	push cx			; cx -> stack

	mov ax, cx			; ax = cx
	;--- call the "writenum" function ---
	call writenum

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, bottle		; parameter #2 is &bottle
	mov rdx, _length_of_bottle		; parameter #3 is len(bottle)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, plural		; parameter #2 is &plural
	mov rdx, _length_of_plural		; parameter #3 is len(plural)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, ofbeer		; parameter #2 is &ofbeer
	mov rdx, _length_of_ofbeer		; parameter #3 is len(ofbeer)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, wall			; parameter #2 is &wall
	mov rdx, _length_of_wall		; parameter #3 is len(wall)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, sep			; parameter #2 is &sep
	mov rdx, _length_of_sep		; parameter #3 is len(sep)
	syscall				; perform the call

	pop ax				; stack -> ax

	;--- call the "writenum" function ---
	call writenum

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, bottle		; parameter #2 is &bottle
	mov rdx, _length_of_bottle		; parameter #3 is len(bottle)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, plural		; parameter #2 is &plural
	mov rdx, _length_of_plural		; parameter #3 is len(plural)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, ofbeer		; parameter #2 is &ofbeer
	mov rdx, _length_of_ofbeer		; parameter #3 is len(ofbeer)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, dotnl			; parameter #2 is &dotnl
	mov rdx, _length_of_dotnl		; parameter #3 is len(dotnl)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, takedown		; parameter #2 is &takedown
	mov rdx, _length_of_takedown		; parameter #3 is len(takedown)
	syscall				; perform the call

	pop ax				; stack -> ax

	dec ax			; ax--
	push ax			; ax -> stack

	;--- call the "writenum" function ---
	call writenum

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, bottle		; parameter #2 is &bottle
	mov rdx, _length_of_bottle		; parameter #3 is len(bottle)
	syscall				; perform the call

	pop ax				; stack -> ax

	;--- call the "pluralifneeded" function ---
	call pluralifneeded

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, ofbeer		; parameter #2 is &ofbeer
	mov rdx, _length_of_ofbeer		; parameter #3 is len(ofbeer)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, wall			; parameter #2 is &wall
	mov rdx, _length_of_wall		; parameter #3 is len(wall)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, dotnl			; parameter #2 is &dotnl
	mov rdx, _length_of_dotnl		; parameter #3 is len(dotnl)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, nl			; parameter #2 is &nl
	mov rdx, _length_of_nl		; parameter #3 is len(nl)
	syscall				; perform the call

	pop rcx				; restore counter
	dec rcx				; decrease counter
	jz l2_end			; jump out if the loop is done
	cmp cx, 2			; compare
	jge l2			; continue

	mov ax, 1		; ax = 1
	;--- call the "writenum" function ---
	call writenum

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, bottle		; parameter #2 is &bottle
	mov rdx, _length_of_bottle		; parameter #3 is len(bottle)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, ofbeer		; parameter #2 is &ofbeer
	mov rdx, _length_of_ofbeer		; parameter #3 is len(ofbeer)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, wall			; parameter #2 is &wall
	mov rdx, _length_of_wall		; parameter #3 is len(wall)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, sep			; parameter #2 is &sep
	mov rdx, _length_of_sep		; parameter #3 is len(sep)
	syscall				; perform the call

	mov ax, 1		; ax = 1
	;--- call the "writenum" function ---
	call writenum

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, bottle		; parameter #2 is &bottle
	mov rdx, _length_of_bottle		; parameter #3 is len(bottle)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, ofbeer		; parameter #2 is &ofbeer
	mov rdx, _length_of_ofbeer		; parameter #3 is len(ofbeer)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, dotnl			; parameter #2 is &dotnl
	mov rdx, _length_of_dotnl		; parameter #3 is len(dotnl)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, takedown		; parameter #2 is &takedown
	mov rdx, _length_of_takedown		; parameter #3 is len(takedown)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, l_no			; parameter #2 is &l_no
	mov rdx, _length_of_l_no		; parameter #3 is len(l_no)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, more			; parameter #2 is &more
	mov rdx, _length_of_more		; parameter #3 is len(more)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, wall			; parameter #2 is &wall
	mov rdx, _length_of_wall		; parameter #3 is len(wall)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, dotnl			; parameter #2 is &dotnl
	mov rdx, _length_of_dotnl		; parameter #3 is len(dotnl)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, nl			; parameter #2 is &nl
	mov rdx, _length_of_nl		; parameter #3 is len(nl)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, u_no			; parameter #2 is &u_no
	mov rdx, _length_of_u_no		; parameter #3 is len(u_no)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, more			; parameter #2 is &more
	mov rdx, _length_of_more		; parameter #3 is len(more)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, wall			; parameter #2 is &wall
	mov rdx, _length_of_wall		; parameter #3 is len(wall)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, sep			; parameter #2 is &sep
	mov rdx, _length_of_sep		; parameter #3 is len(sep)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, l_no			; parameter #2 is &l_no
	mov rdx, _length_of_l_no		; parameter #3 is len(l_no)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, more			; parameter #2 is &more
	mov rdx, _length_of_more		; parameter #3 is len(more)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, dotnl			; parameter #2 is &dotnl
	mov rdx, _length_of_dotnl		; parameter #3 is len(dotnl)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, store			; parameter #2 is &store
	mov rdx, _length_of_store		; parameter #3 is len(store)
	syscall				; perform the call

	mov ax, 99		; ax = 99
	;--- call the "writenum" function ---
	call writenum

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, bottle		; parameter #2 is &bottle
	mov rdx, _length_of_bottle		; parameter #3 is len(bottle)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, plural		; parameter #2 is &plural
	mov rdx, _length_of_plural		; parameter #3 is len(plural)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, ofbeer		; parameter #2 is &ofbeer
	mov rdx, _length_of_ofbeer		; parameter #3 is len(ofbeer)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, wall			; parameter #2 is &wall
	mov rdx, _length_of_wall		; parameter #3 is len(wall)
	syscall				; perform the call

	;--- system call ---
	mov rax, 1			; function call: 1
	mov rdi, 1			; parameter #1 is 1
	mov rsi, dotnl			; parameter #2 is &dotnl
	mov rdx, _length_of_dotnl		; parameter #3 is len(dotnl)
	syscall				; perform the call

	pop rcx				; restore counter
	dec rcx				; decrease counter
	jnz l2				; loop until rcx is zero
l2_end:				; end of loop l2
	;--- end of loop l2 ---


	;--- return from "main" ---
	mov rax, 60			; function call: 60
	xor rdi, rdi			; return code 0
	syscall				; exit program



