/* 
 * "sprinter" uses C funtion calling conventions, while "myitoa" does not.
 * sprinter is a simplified implementation of the C sprintf function, while myitoa is a helper function for converting an integer to a string.
 * The c file is a draft for this file, which is why the assembly is so "C-like" (it was an experiment).
 */

.data

/*** Variables for the "myitoa" function ***/

n:
    .int 0

ptr_s:
    .int 0

b:
    .int 0

digits:
    /* static char digits[] = "0123456789abcdefghijklmnopqrstuvwxyz"; */
    .string "0123456789abcdefghijklmnopqrstuvwxyz\0"

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

ptr_p1:
    /* char *p1; */
    .int 0

ptr_p2:
    /* char *p2; */
    .int 0

reverse:
    .rept 4096	    /* max length */
    .string "\0"
    .endr


/********** Variables for the sprinter function ***********/

/* For finding the % codes */

lettercursor:
    .int 0

inCode:
    .int -1	# -1 is false

letter:
    .byte 0

lenRetString:
    .int 0

codeSoFar:
    /* Assuming the formatting codes won't be longer than this */
    .rept 4096		/* Max length */
    .string "\0"
    .endr

lenCodeSoFar:
    .int 0

formattingCodeLetter:
    # Formatting letter ('c', 'd', 's' or 'x')
    .byte '\0'

/* Used when converting from stringnumber to integer */

tenth:
    .int 0

charcounter:
    .int 0

minWidth:
    .int 0

digit:
    .int 0	# byte is too small for the conversion calculations

expo:
    .int 0

/* Used when getting arguments that were given to the function *//

chararg:
    .byte 0

intarg:
    .int 0

numberstring:
    /* Assuming the digits for converted numbers won't be longer than this */
    .rept 99		    /* Max length */
    .string "\0"
    .endr

ptr_strarg:
    /* char* strarg */
    .int 0

strlen:
    .int 0

ptr_res:
    /* char* res */
    .int 0

ptr_format:
    /* char* format */
    .int 0

argbytecounter:
    /* Used for keeping track of the next var args */
    .int 0

.text

    /* int sprinter(char *res, char *format, ...) */
    .globl sprinter
sprinter:

    pushl   %ebp
    movl    %esp, %ebp

    # The program segfauls if I push these and pop them in the opposite order later on
    #pushl   %ebx
    #pushl   %edi
    #pushl   %esi

    movl    $0, %ecx

    /* Blank variables */

    movl    $0, lettercursor
    movl    $-1, inCode
    movb    $0, (letter)
    movl    $0, lenRetString
    movb    $0, (codeSoFar)
    movl    $0, lenCodeSoFar
    movb    $0, (formattingCodeLetter)
    movl    $0, tenth
    movl    $0, charcounter
    movl    $0, minWidth
    movl    $0, digit
    movl    $0, expo
    movb    $0, chararg
    movl    $0, intarg
    movb    $0, (numberstring)
    movl    $0, ptr_strarg
    movl    $0, ptr_format
    movl    $0, argbytecounter

    /* Load the parameters */

    movl    8(%esp), %eax
    movl    %eax, ptr_res

    movl    12(%esp), %eax
    movl    %eax, ptr_format

    /* The rest of the arguments should be fetched from 16(%esp) etc */
    movl    $16, argbytecounter

    # letter = format[lettercursor];
    movl    lettercursor, %eax
    addl    ptr_format, %eax
    xorl    %ebx, %ebx		/* zero out ebx first */
    movb    (%eax), %bl
    movb    %bl, letter

move_to_next_letter:
    /* inCode == 0; // true? */
    cmpl    $0, inCode
    je	    add_letter_to_code_so_far

    jmp just_a_regular_letter

add_letter_to_code_so_far:

    # Add a letter to the codeSoFar string if it's not \0

    /* letter == 0; */
    cmpb    $0, letter
    je	    over2

    /* Gather code letters */
    /* codeSoFar[lenCodeSoFar] = letter; */
    movl    $codeSoFar, %eax
    addl    lenCodeSoFar, %eax
    xorl    %ebx, %ebx		/* zero out ebx first */
    movb    letter, %bl
    movb    %bl, (%eax)

    /* lenCodeSoFar++; */
    incl    lenCodeSoFar
over2:

    /* See if the letter is 0, % or a letter from a/A to z/Z */

    /* The checks should work like this:
     * 0 : letter ok
     * % : letter ok
     * g : check_in_az -> letter ok
     * G : check_in_AZ -> letter ok
     * ! : check_in_az -> check_over_A -> check_in_AZ -> not ok
     */

    xorl    %ebx, %ebx		/* zero out ebx first */
    movb    letter, %bl

    cmpb    $0, %bl		# check if the letter is \0
    je	    ok_letter

    cmpb    $'%', %bl		# check if the letter is %
    je	    ok_letter

    cmpb    $'a', %bl		# check if the letter is >='a'
    jge	    check_in_az

check_over_A:
    cmpb    $'A', %bl		# check if the letter is >='A'
    jge	    check_in_AZ
    jmp	    letter_check_done	/* letter is not ok */

check_in_az:
    cmpb    $'z', %bl		# check if the letter is <='z'
    jle	    ok_letter
    jmp	    check_over_A

check_in_AZ:
    cmpb    $'Z', %bl		# check if the letter is <='Z'
    jle	    ok_letter
    jmp	    letter_check_done	/* letter is not ok */

ok_letter:

    /* interpret the code (after the first %) */

    /* ((lenCodeSoFar == 1) && (codeSoFar[0] == '%')); */
    cmpl    $1, lenCodeSoFar
    jne	    check0_fail

    /* check0 ok */
    movl    $codeSoFar, %eax
    movb    (%eax), %bl
    cmpb    $'%', %bl
    jne	    check0_fail

    /* both ok */
    jmp	    just_percentage

check0_fail:

    # check: lenCodeSoFar > 0;
    cmpl    $0, lenCodeSoFar
    jg	    more_formatting_code
    jmp	    no_formatting_code

just_percentage:

    # Just a % sign
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    $'%', (%eax)

    incl    lenRetString
    jmp	    interpret_letter_done

more_formatting_code:

    # eax = lenCodeSoFar - 1;
    movl    lenCodeSoFar, %eax
    decl    %eax

    # formattingCodeLetter = codeSoFar[eax];
    xorl    %ebx, %ebx		/* zero out ebx first */
    addl    $codeSoFar, %eax
    movb    (%eax), %bl
    movb    %bl, formattingCodeLetter

    /*
     * Find the minimum width (number before the formatting letter)
     */

    # minWidth = 0;
    movl    $0, minWidth
    # charcounter = 0;
    movl    $0, charcounter

next_digit:

    # check: charcounter >= (lenCodeSoFar - 1);
    movl    lenCodeSoFar, %eax
    decl    %eax
    cmpl    %eax, charcounter
    jge	    check_formatting_letter

    /* Substract 2 to get the right N for 10^N when converting the
     * number from a string to an int.
     */

    # eax = lenCodeSoFar - 2;
    decl    %eax

    # tenth = eax - charcounter;
    subl    charcounter, %eax
    movl    %eax, tenth

    # Check if it's a digit between '0' and '9'
    movl    $codeSoFar, %eax
    addl    charcounter, %eax
    xorl    %ebx, %ebx			/* zero out ebx first */
    movb    (%eax), %bl

    cmpb    $'0', %bl
    jl	    not_a_digit

    cmpb    $'9', %bl
    jg	    not_a_digit

is_a_digit:

    # Convert char digit to numeric digit by subtracting '0'
    subb    $'0', %bl
    movb    $0, %bh
    movl    %ebx, digit

    /*
     * Make digits into a number
     * using digit * 10^tenth
     */

    movl    $0, expo
process_digit:
    # check: expo >= tenth;
    movl    expo, %eax
    cmpl    tenth, %eax
    jge	    stop_processing_digit

    # digit *= 10
    movl    digit, %eax
    shll    $1, %eax
    movl    %eax, %ebx
    shll    $1, %eax
    shll    $1, %eax
    addl    %ebx, %eax
    movl    %eax, digit

    # expo++
    incl    expo
    jmp	    process_digit

stop_processing_digit:

    # minWidth += digit;
    movl    minWidth, %eax
    addl    digit, %eax
    movl    %eax, minWidth

    jmp	    digit_check_done

not_a_digit:

    jmp	    return_error

digit_check_done:

    # charcounter++;
    incl    charcounter
    jmp	    next_digit

check_formatting_letter:

    /* Check what the formatting letter is */

    cmpb    $'c', formattingCodeLetter
    je	    letter_is_c
    cmpb    $'d', formattingCodeLetter
    je	    letter_is_d
    cmpb    $'s', formattingCodeLetter
    je	    letter_is_s
    cmpb    $'x', formattingCodeLetter
    je	    letter_is_x

    jmp	    letter_is_unknown


/******************* Formatting letter is 'c' *************************/

letter_is_c:

    /* Get the next arguments
     * must use int instead of char because of how va_arg works in C
     */
    # chararg = (char)va_arg(args, int);
    movl    %esp, %ebx
    addl    argbytecounter, %ebx
    movl    (%ebx), %eax
    addl    $4, argbytecounter
    movb    %al, chararg

    # Add minWidth - 1 spaces in front of the char

    # ecx = 0
    xorl    %ecx, %ecx
add_next_char:
    # eax = minWidth - 1;
    movl    minWidth, %eax
    decl    %eax

    cmpl    %eax, %ecx
    jge	    stop_adding_chars

    # ptr_res[lenRetString] = ' ';
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    $' ', (%eax)

    # lenRetString++;
    incl    lenRetString

    # ecx++;
    incl    %ecx

    jmp	    add_next_char
stop_adding_chars:
    /* Add char argument to the ptr_res */

    # ptr_res[lenRetString] = chararg;
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    xorl    %ebx, %ebx			/* zero out ebx first */
    movb    chararg, %bl
    movb    %bl, (%eax)

    # lenRetString++;
    incl    lenRetString

    jmp	    letter_done

/******************* Formatting letter is 'd' *************************/

letter_is_d:
    /* Get the next arguments
     * must use int instead of char because of how va_arg works in C
     */
    # intarg = va_arg(args, int);
    movl    %esp, %ebx
    addl    argbytecounter, %ebx
    movl    (%ebx), %eax
    addl    $4, argbytecounter
    movl    %eax, intarg

    /*
     * Build up a string based on the decimal number in intarg
     */

    movl    intarg, %eax
    movl    %eax, n
    movl    $numberstring, %eax
    movl    %eax, ptr_s
    movl    $10, b
    call    myitoa

    /* Find the length of numberstring */

    # ecx = 0
    xorl    %ecx, %ecx

check_next_character:
    # check: numberstring[i] == '\0';
    movb    numberstring(%ecx), %al
    cmpb    $0, %al
    je	    stop_checking_next_character

    # i++;
    incl    %ecx

    jmp check_next_character
stop_checking_next_character:

    # strlen = i
    movl    %ecx, strlen

    /* Add (minWidth-strlen) spaces in front of the number string */

    # i = 0
    xorl    %ecx, %ecx
add_another_character:

    # eax = minWidth - strlen
    movl    minWidth, %eax
    subl    strlen, %eax

    # check: i >= eax
    cmpl    %eax, %ecx
    jge	    stop_adding_characters
    
    # ptr_res[lenRetString] = ' ';
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    $' ', (%eax)

    # lenRetString++;
    incl    lenRetString

    # i++;
    incl    %ecx

    jmp	    add_another_character
stop_adding_characters:

    /* Add number string to the ptr_res */

    # ecx = 0
    xorl    %ecx, %ecx
add_digit:
    # check: ecx >= strlen;
    cmpl    strlen, %ecx
    jge	    stop_adding_digits

    # ptr_res[lenRetString] = numberstring[i];
    xorl    %ebx, %ebx			/* zero out ebx first */
    movb    numberstring(%ecx), %bl
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    %bl, (%eax)

    # lenRetString++;
    incl    lenRetString

    # ecx++;
    incl    %ecx

    jmp	    add_digit

stop_adding_digits:

    jmp	    letter_done


/******************* Formatting letter is 's' *************************/

letter_is_s:

    /* Get the next arguments */

    # ptr_strarg = (char*)va_arg(args, char*);
    movl    %esp, %ebx
    addl    argbytecounter, %ebx
    movl    (%ebx), %eax
    addl    $4, argbytecounter
    movl    %eax, ptr_strarg

    # strlen=0; (using ecx for strlen)
    xorl    %ecx, %ecx
check_next_character2:
    #check: ptr_strarg[strlen] == 0; (using ecx for strlen)
    movl    ptr_strarg, %eax
    addl    %ecx, %eax

    movb    (%eax), %bl
    cmpb    $0, %bl
    je	    stop_checking_next_character2

    # strlen++; (using ecx for strlen)
    incl    %ecx

    jmp	    check_next_character2
stop_checking_next_character2:
    # strlen = ecx
    movl    %ecx, strlen

    /* Add (minWidth-strlen) spaces in front of the string */

    # ecx = 0
    xorl    %ecx, %ecx
add_another_space:
    # eax = minWidth - strlen;
    movl    minWidth, %eax
    subl    strlen, %eax

    # check: ecx >= eax;
    cmpl    %eax, %ecx
    jge	    stop_adding_spaces

    # ptr_res[lenRetString] = ' ';
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    $' ', (%eax)

    # lenRetString++;
    incl    lenRetString

    #ecx++;
    incl    %ecx

    jmp	    add_another_space
stop_adding_spaces:

    /* Add string argument to the ptr_res */

    # i = 0
    xorl    %ecx, %ecx
copy_another_character:
    # check: i >= strlen;
    cmpl    strlen, %ecx
    jge	    stop_copying_characters

    # ptr_res[lenRetString] = strarg[i];
    movl    ptr_strarg, %eax
    addl    %ecx, %eax
    xorl    %ebx, %ebx			/* zero out ebx first */
    movb    (%eax), %bl
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    %bl, (%eax)

    # lenRetString++;
    incl    lenRetString

    # i++;
    incl    %ecx

    jmp	    copy_another_character

stop_copying_characters:

    jmp	    letter_done;

/******************* Formatting letter is 'x' *************************/

letter_is_x:

    /* Get the next arguments
     * must use int instead of char because of how va_arg works in C
     */
    # intarg = va_arg(args, int);
    movl    %esp, %ebx
    addl    argbytecounter, %ebx
    movl    (%ebx), %eax
    addl    $4, argbytecounter
    movl    %eax, intarg

    /*
     * Build up a string based on the decimal number in intarg
     */

    movl    intarg, %eax
    movl    %eax, n
    movl    $numberstring, %eax
    movl    %eax, ptr_s
    movl    $16, b	    /* hex */
    call    myitoa

    /* Find the length of numberstring */

    # i = 0;
    xorl    %ecx, %ecx
check_next_character3:
    # check: numberstring[i] == '\0';
    movb    numberstring(%ecx), %al
    cmpb    $0, %al
    je	    stop_checking_next_character3

    # i++;
    incl    %ecx

    jmp	    check_next_character3
stop_checking_next_character3:

    # strlen = i;
    movl    %ecx, strlen

    /* Add (minWidth-strlen) spaces in front of the number string */

    # i = 0;
    xorl    %ecx, %ecx
add_another_space2:
    # eax = minWidth - strlen;
    movl    minWidth, %eax
    subl    strlen, %eax

    # check: i >= eax;
    cmpl    %eax, %ecx
    jge	    stop_adding_another_space2

    # ptr_res[lenRetString] = ' ';
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    $' ', (%eax)

    # lenRetString++;
    incl    lenRetString

    # i++;
    incl    %ecx

    jmp	    add_another_space2
stop_adding_another_space2:

    /* Add number string to the ptr_res */

    # i = 0;
    xorl    %ecx, %ecx

copy_another_character2:
    # check: i >= strlen;
    cmpl    strlen, %ecx
    jge	    stop_copying_another_character2

    # ptr_res[lenRetString] = numberstring[i];
    xorl    %ebx, %ebx			/* zero out ebx first */
    movb    numberstring(%ecx), %bl
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    movb    %bl, (%eax)

    # lenRetString++;
    incl    lenRetString

    # ecx++;
    incl    %ecx

    jmp	    copy_another_character2
stop_copying_another_character2:

    jmp	    letter_done
letter_is_unknown:
    jmp	    return_error
letter_done:
    jmp	    interpret_letter_done
no_formatting_code:
    jmp	    return_error
interpret_letter_done:

/************** Done interpreting a formatting string ***************************/

    # inCode = -1; // false
    movl    $-1, inCode

    jmp	    letter_check_done

just_a_regular_letter:
    cmpb    $'%', letter
    je	    letter_is_percentage
    jmp	    letter_is_not_percentage

letter_is_percentage:
    /* Entering formatting code */
    # inCode = 0; // true
    movl    $0, inCode

    # lenCodeSoFar = 0;
    movl    $0, lenCodeSoFar

    jmp	    letter_check_done

letter_is_not_percentage:
    cmpb    $0, letter
    je	    letter_check_done

    # ptr_res[lenRetString] = letter;
    movl    ptr_res, %eax
    addl    lenRetString, %eax
    xorl    %ebx, %ebx			/* zero out ebx first */
    movb    letter, %bl
    movb    %bl, (%eax)

    # lenRetString++;
    incl    lenRetString

    # ptr_res[lenRetString] = 0;
    incl    %eax     /* same number as 
		        movl    ptr_res, %eax
			addl    lenRetString, %eax
		      */
    movb    $0, (%eax)	    # set ptr_res[lenRetString] to 0

letter_check_done:

    /* Stop at end of string */

    #check: letter == 0;
    cmpb    $0, letter
    je	    stop_moving_to_next_letter

    /* Next letter of res */

    # lettercursor++;
    incl    lettercursor

    # letter = format[lettercursor];
    movl    ptr_format, %eax
    addl    lettercursor, %eax
    xorl    %ebx, %ebx			/* zero out ebx first */
    movb    (%eax), %bl
    movb    %bl, letter

    jmp	    move_to_next_letter
stop_moving_to_next_letter:

return0:

    # prepare to return the length of the resulting string
    movl    lenRetString, %eax

    # res[lenRetString] = '\0'
    addl    ptr_res, %eax
    movb    $0, (%eax)
    subl    ptr_res, %eax

return:

    # I had no luck enabling these
    #popl   %esi
    #popl   %edi
    #popl   %ebx

    movl    %ebp, %esp
    popl    %ebp

    ret				    # return from the sprinter function

return_error:
    movl    $-1, %eax

/* itoa function
 *
 * converts a number to a string
 *
 * input:
 *	n : int, 
 *	ptr_s : char*, pointer to string to convert
 *	b : int, base number (10 for dec, 16 for hex)
 *
 * output:
 *	puts a char* pointer in %eax with the converted number
 */
myitoa:

    /* Blank variables */

    movl    $0, x
    movl    $0, negative
    movl    $0, sign
    movl    $0, stringlength
    movb    $0, (reverse)
    
    # See if the number is 0
    cmpl    $0, n
    jne	    itoa_over00

    # It's 0, return "0\0"

    movl    ptr_s, %edi
    movb    $'0', (%edi)
    incl    %edi
    movb    $0, (%edi)
    jmp	    itoa_return0

itoa_over00:

    # Check if the number is not the special number 0x80000000

    # n != x
    movl    n, %eax
    cmpl    x, %eax
    jne	    itoa_over0

    /* If it's the special number 0x80000000, copy over the letters from
     * the special string -2147483648 and return that
     */

    # i = 0 (using %edx)
    xorl    %edx, %edx
itoa_loop0:
    
    /* Check if end of string is reached */

    # special[i] == '\0';
    cmpb    $0, special(%edx)
    je itoa_out0

    /* Copy over characters */
    
    # s[i] = special[i];
    pushl   %esi
    pushl   %edi
    movl    $special, %esi	/* source adr in %esi */
    movl    ptr_s, %edi		/* destination adr in %edi */
    movl    %edx, %ecx		/* bytes to move in %ecx */
    rep	    movsb		/* copy bytes over */
    popl    %edi
    popl    %esi
    
    incl    %edx
    jmp	    itoa_loop0

itoa_out0:

    # return s
    jmp	    itoa_return0

itoa_over0:

    /* Make n positive and store the sign */

    # sign = n;
    movl    n, %eax
    movl    %eax, sign

    /* skip if positive */

    # sign >= 0;
    movl    sign, %eax
    cmp	    $0, %eax
    jge	    itoa_positive

    /* If it's a negative hex */

    cmpl    $16, b	# is base 16?
    jne	    not_base_16

    movl    $28, %ebx	# pos = 28
    xorl    %ecx, %ecx	# i = 0

neghexloop:

    /* s[i] = digits[((n >> pos) & (b-1))]; */

    # %ecx is now "temp", not "i"
    pushl   %ecx

    # move "b - 1" into temp (%ecx)
    movl    b, %ecx
    decl    %ecx

    # move "n >> pos" into %eax
    movl    n,	%eax
    pushl   %ecx
    xorl    %ecx, %ecx	    # %ecx = 0
    mov	    %bl, %cl
    shr     %cl, %eax	    # seems like I have to use %cl, hence the push&pop of %ecx
    popl    %ecx

    # and %ecx ("temp") with %eax (the result of n >> pos)
    andl    %ecx, %eax

    # %ecx is now back to "i", not "temp"
    popl    %ecx

    # digit position should now be in %eax
    
    # get the digit from the digits variable
    movl    $digits, %esi
    addl    %eax, %esi

    # the digit should now be at %esi, move it over to (ptr_s + %ecx)
    movb    (%esi), %al
    movl    ptr_s, %edi
    addl    %ecx, %edi

    # save %ecx when using movsb to move N bytes where N is %ecx
    pushl   %ecx
    movl    $1, %ecx
    movsb
    popl    %ecx

    # continue to the next hex digit

    subl    $4, %ebx
    incl    %ecx
    cmpl    $0, %ebx
    jl	    done_neg_hex_digits

    jmp	    neghexloop
done_neg_hex_digits:
    # s[i] = '\0';	# i is ecx
    movl    ptr_s, %eax
    addl    %ecx, %eax
    movb    $0, (%eax)
    jmp	    itoa_return0

not_base_16:

    # n = -n;
    negl    n

itoa_positive:

    /* Add the digits to the string */

    # i = 0; (using ecx)
    xorl    %ecx, %ecx

itoa_loop1:
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
    pushl   %ecx
    pushl   %edi
    pushl   %esi

    movl    %ebx, %eax	    /* Store ebx */
    movl    $digits, %ebx   /* The adr of digits in ebx */
    addl    %edx, %ebx	    /* Add edx to get the right digit */
    movl    %ebx, %esi	    /* Put the address in esi */
    movl    %eax, %ebx	    /* Restore ebx */
    movl    ptr_s, %edi	    /* The s pointer in edi */
    addl    %ecx, %edi	    /* Add i to get the right char */

    pushl   %ecx
    movl    $1, %ecx	    /* Number of bytes to copy */
    movsb		    /* Copy it over */
    popl    %ecx

    popl    %esi
    popl    %edi
    popl    %ecx

    # i++;
    incl    %ecx

    # n = result0; (ebx)
    movl    %ebx, %eax
    movl    %eax, n

    # result > 0; (ebx)
    cmpl    $0, %ebx
    jg	    itoa_loop1

    /* Add '-' if needed */

    # sign >= 0;
    cmpl    $0, sign
    jg	    itoa_over1

    # s[i] = '-' (using ecx for i)
    movl    ptr_s, %eax
    addl    %ecx, %eax
    movb    $'-', %bl
    movb    %bl, (%eax)

    # i++
    incl    %ecx

itoa_over1:
    # s[i] = '\0'; (using ecx for i)
    movl    ptr_s, %eax
    addl    %ecx, %eax
    movb    $0, (%eax)

    /* If there's no string, return that */

    # !s || !*s;
    cmpl    $0, ptr_s
    jne	    itoa_skip1
    movb    (ptr_s), %al
    cmpb    $0, %al
    jne	    itoa_skip1

    # return s
    jmp	    itoa_return0

itoa_skip1:
    /* Find the length of ptr_s */

    xorl    %ecx, %ecx
length_loop:
    movl    ptr_s, %eax
    addl    %ecx, %eax
    movb    (%eax), %bl
    cmp	    $0, %bl
    je	    length_loop_out
    incl    %ecx
    jmp length_loop
length_loop_out:

    /* Length (not including \0) is now in ecx */

    /* Save length in edx */
    movl    %ecx, %edx

    cmpl    $1, %ecx
    jle	    itoa_return0

    /* Try to reverse the number */
    
    # p1 = p
    movl    ptr_s, %eax
    movl    %eax, ptr_p1
    # p2 = p + strlen(p) - 1
    movl    ptr_s, %eax
    addl    %ecx, %eax
    decl    %eax
    movl    %eax, ptr_p2

    /* Loop for reversing the number before returning */

itoa_revloop:
    # check: p1 < p2
    movl    ptr_p1, %eax
    movl    ptr_p2, %ebx
    cmpl    %ebx, %eax
    jge	    itoa_revloop_out

    # c = *p1
    movl    ptr_p1, %eax
    movb    (%eax), %dl
    # *p1 = *p2
    movl    ptr_p2, %ebx
    movb    (%ebx), %al
    movl    ptr_p1, %ebx
    movb    %al, (%ebx)
    # *p2 = c
    movl    ptr_p2, %ebx
    movb    %dl, (%ebx)
    # p1++
    incl    ptr_p1
    # p2--
    decl    ptr_p2

    jmp itoa_revloop
itoa_revloop_out: 


itoa_return0:
    # return char* ptr_s
    movl    ptr_s, %eax

    ret					# return from the myitoa function

