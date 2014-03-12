.data

n:
    .int 0

ptr_s:
    .int 0

b:
    .int 0

digits:
    /* static char digits[] = "0123456789abcdefghijklmnopqrstuvwxyz"; */
    .string "0123456789abcdefghijklmnopqrstuvwxyz\0"

reverse:
    .string "\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0"

x:
    /* int x = 0x80000000; */
    .int 0x80000000

special:
    /* char* special = "-2147483648"; */
    .string "-2147483648\0"

negative:
    /* int negative; */
    .int 0

sign:
    /* int sign; */
    .int 0

stringlength:
    /* int stringlength; */
    .int 0

ptr_1:
    /* char *p1; */
    .int 0

ptr_2:
    /* char *p2; */
    .int 0


.text

    /* char *myitoa(int n, char *s, int b) */
    .globl myitoa
myitoa:

    pushl   %ebp
    movl    %esp, %ebp

    movl    $0, %ecx

    /*
    %esp has the return address, 4(%esp) is the first param,
    8(%esp) is the next param etc., () fetches the value at
    whatever is within. N() fetches the value at whatever is
    within + N 32-bit instructions typically fetch 32-bit values,
    which makes 4(N), 8(N), 12(N) etc usual.
    */

    /* Load the parameters */

    movl    8(%esp), %eax
    movl    %eax, n

    movl    12(%esp), %eax
    movl    %eax, ptr_s

    movl    16(%esp), %eax
    movl    %eax, b

    # Check if the number is not the special number 0x80000000

    # n != x
    movl    n, %eax
    cmpl    x, %eax
    jne over0

    /* If it's the special number 0x80000000, copy over the letters from
     * the special string -2147483648 and return that
     */

    # i = 0 (using %edx)
    xorl    %edx, %edx
loop0:
    
    /* Check if end of string is reached */

    # special[i] == '\0';
    cmpb    $0, special(%edx)
    je out0

    /* Copy over characters */
    
    # s[i] = special[i];
    movl    $special, %esi	/* source adr in %esi */
    movl    ptr_s, %edi		/* destination adr in %edi */
    movl    %edx, %ecx		/* bytes to move in %ecx */
    rep	    movsb		/* copy bytes over */
    
    incl    %edx
    jmp loop0
out0:

    # return s
    jmp	    return0

over0:

    /* Make n positive and store the sign */

    # sign = n;
    movl    n, %eax
    movl    %eax, sign

    /* skip if positive */

    # sign >= 0;
    movl    sign, %eax
    cmp	    $0, %eax
    jge	    positive

    # n = -n;
    negl    n

positive:

    /* Add the digits to the string */

    # i = 0; (using ecx)
    xorl    %ecx, %ecx

loop1:
    /* Trying to do this: s[i] = digits[n % b]; */

    movl    n, %eax
    movl    b, %ebx

    # n % b
    xorl    %edx, %edx
    divl    %ebx

    /* n / b is now in eax and
       n % b is now in edx
     */

    # result0 = n / b; (use ebx as result0)
    movl   %eax, %ebx
    	
    /* s[i] = digits[n % b]
     *
     * digit pos is now in %dl
     * digits is a variable
     * ptr_s is the destination
     * ecx is i
     */
    movl    %ebx, %eax	    /* Store ebx */
    movl    $digits, %ebx   /* The adr of digits in ebx */
    addl    %edx, %ebx	    /* Add edx to get the right digit */
    movl    %ebx, %esi	    /* Put the address in esi */
    movl    %eax, %ebx	    /* Restore ebx */
    movl    ptr_s, %edi	    /* The s pointer in edi */
    addl    %ecx, %edi	    /* Add i to get the right char */
    movl    %ecx, %eax	    /* Store ecx */
    movl    $1, %ecx	    /* Number of bytes to copy */
    movsb		    /* Copy it over */
    movl    %eax, %ecx	    /* Restore ecx */

    # i++;
    incl    %ecx

    # n = result0; (ebx)
    movl    %ebx, %eax
    movl    %eax, n

    # result > 0; (ebx)
    cmpl    $0, %ebx
    jg	    loop1

    /* Add '-' if needed */

    # sign >= 0;
    cmpl    $0, sign
    jg	    over1

    # s[i] = '-' (using ecx for i)
    movl    ptr_s, %eax
    addl    %ecx, %eax
    movb    $'-', %bl
    movb    %bl, (%eax)

    # i++
    incl    %ecx

over1:
    # s[i] = '\0'; (using ecx for i)
    movl    ptr_s, %eax
    addl    %ecx, %eax
    movb    $0, (%eax)

    /* If there's no string, return that */

    # !s || !*s;
    cmpl    $0, ptr_s
    jne	    skip1
    movb    (ptr_s), %al
    cmpb    $0, %al
    jne	    skip1

    # return s
    jmp	    return0

skip1:

    /* Find the length of the string */

    # i = 0; (using edx for i)
    xorl    %edx, %edx

loop2:
    # i++
    incl    %edx

    /* check: s[i] != '\0'; */
    movl    ptr_s, %eax
    addl    %edx, %eax
    cmpb    $0, (%eax)
    jne	    loop2
    
    /* Reverse the string */
    movl    ptr_s, %eax	    # eax = ptr_s + len(ptr_s)
    addl    %edx, %eax
    decl    %eax	    # don't copy over the \0
    movl    $reverse, %ebx  # ebx = reverse
    movl    %edx, %ecx	    # counter

rev_loop:
    movb    (%eax), %dl
    movb    %dl, (%ebx)
    decl    %eax
    incl    %ebx
    decl    %ecx

    cmp	    $0, %ecx
    jg	    rev_loop

    /* Add 0 at the end of the reverse string, just to be sure */
    movb    $0, (%ebx)

    /* copy the reversed string over to ptr_s */
    movl    $reverse, %eax
    movl    ptr_s, %ebx

copy_loop:
    movb    (%eax), %dl

    # if bl == 0, jump out
    cmpb    $0, %dl
    je	    copy_done

    movb    %dl, (%ebx)

    incl    %eax
    incl    %ebx

    jmp copy_loop
copy_done:


return0:
    # return char* ptr_s
    movl    ptr_s, %eax
    popl    %ebp
    ret

