; Generated with Battlestar 0.3, at 2014-10-29 11:17

bits 32

; Thanks to http://wiki.osdev.org/Bare_Bones_with_NASM

; Declare constants used for creating a multiboot header.
MBALIGN     equ  1<<0                   ; align loaded modules on page boundaries
MEMINFO     equ  1<<1                   ; provide memory map
FLAGS       equ  MBALIGN | MEMINFO      ; this is the Multiboot 'flag' field
MAGIC       equ  0x1BADB002             ; 'magic number' lets bootloader find the header
CHECKSUM    equ -(MAGIC + FLAGS)        ; checksum of above, to prove we are multiboot
 
; Declare a header as in the Multiboot Standard. We put this into a special
; section so we can force the header to be in the start of the final program.
; You don't need to understand all these details as it is just magic values that
; is documented in the multiboot standard. The bootloader will search for this
; magic sequence and recognize us as a multiboot kernel.
section .multiboot
align 4
	dd MAGIC
	dd FLAGS
	dd CHECKSUM
 
; Currently the stack pointer register (esp) points at anything and using it may
; cause massive harm. Instead, we'll provide our own stack. We will allocate
; room for a small temporary stack by creating a symbol at the bottom of it,
; then allocating 16384 bytes for it, and finally creating a symbol at the top.
section .bootstrap_stack
align 4
stack_bottom:
times 16384 db 0
stack_top:

section .text

extern kmain			; external symbol

;--- function main ---
global main			; make label available to the linker
global _start			; make label available to the linker
_start:				; starting point of the program
	mov esp, stack_top	; set the esp register to the top of the stack (special case for bootable kernels)
main:				; name of the function


	;--- call the "kmain" function ---
	call kmain

	; --- full stop ---
	cli		; clear interrupts
.hang:
	hlt
	jmp .hang	; loop forever




