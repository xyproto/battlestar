#include <stdarg.h>

// Maximum length of strings that are to be formatted
#define MAXLEN 1023

// Maximum number of digits
#define MAXDIGITS 99

char *myitoa(int n, char *s, int b) {

    // For assembly checks
    int check;
    int result0;

    static char digits[] = "0123456789abcdefghijklmnopqrstuvwxyz";
    int x = 0x80000000;
    char* special = "-2147483648";

    int i;
    int negative;
    int sign;
    int stringlength;
    char *p1;
    char *p2;

    // Check if the number is not the special number 0x80000000

    check = n != x;
    if (check) goto over0;

    // If it's the special number 0x80000000, copy over the letters from
    // the special string -2147483648 and return that

    i = 0;
loop0:
    check = special[i] == '\0';
    if (check) goto out0;
    s[i] = special[i];
    i++;
    goto loop0;
out0:

    return s;
over0:

    // Make n positive and store the sign

    sign = n;

    check = sign >= 0;
    if (check) goto skip0; // skip if positive


    if (b != 16) goto skip2;

    int pos = 28; // %ebx
    i = 0; // %ecx
neghexloop:
    s[i] = digits[((n >> pos) & (b-1))];
    pos -= 4;
    i++;
    if (pos < 0) {
	goto out2;
    }
    goto neghexloop;
out2:
    s[i] = '\0';
    goto out3;
skip2:

    n = -n;
skip0:


    // Add the digits to the string
    i = 0;

loop1:
    s[i] = digits[n % b];
    i++;
    result0 = n / b;
    n = result0;
    check = result0 > 0;
    if (check) goto loop1;

    // Add '-' if needed

    check = sign >= 0;
    if (check) goto over1;
    s[i] = '-';
    i++;
over1:
    s[i] = '\0';

    // If there's no string, return that

    check = !s || !*s;
    if (!check) goto skip1;
    return s;
skip1:

    // Find the length of the string

    i = 0;
loop2:
    i++;
    check = s[i] != '\0';
    if (check) goto loop2;
    i--;

    // Reverse the string

    p1 = s;
    p2 = s + i;
loop3:
    check = p2 <= p1;
    if (check) goto out3;
    *p1 ^= *p2;
    *p2 ^= *p1;
    *p1 ^= *p2;
    p1++;
    p2--;
    goto loop3;
out3:

    return s;
}

int sprinter(char *res, char *format, ...) {

    va_list args;

    // For assembly checks
    int check;
    int result0;
    int temp0;
    int temp1;
    int temp2;

    // 1. Find all instances of % that are not %%
    int lettercursor = 0;
    int inCode = -1; // false

    char letter;
    char retString[MAXLEN + MAXLEN]; // Return string can be MAXLEN*2 long maximum
    int lenRetString = 0; // Length of return string;
    char codeSoFar[MAXLEN]; // Formatting code can be MAXLEN long maximum
    int lenCodeSoFar = 0; // Length of formatting code so far
    char formattingCodeLetter; // Formatting letter ('c', 'd', 's' or 'x')

    // Used when converting from stringnumber to integer
    int tenth = 0;
    int charcounter = 0;
    int minWidth = 0;
    int digit = 0;
    int expo = 0;

    // Used when getting arguments that were given to the function
    char chararg = '\0';

    int intarg = 0;

    char numberstring[MAXDIGITS]; // Maximum number of digits for the converted number
    char* strarg;
    int strlen = 0;
    int i;

    va_start(args, format);

    letter = format[lettercursor];

loop4:
    check = inCode == 0; // is true?

    if (check) goto add_letter_to_code_so_far;
    goto just_a_regular_letter;

add_letter_to_code_so_far:

    // Add a letter to the codeSoFar string if it's not \0

    check = letter == 0;
    if (check) goto over2;
    // Gather code letters
    codeSoFar[lenCodeSoFar] = letter;
    lenCodeSoFar++;
over2:

    // See if the letter is 0, % or a letter from a/A to z/Z
    temp0 = letter == 0;
    temp1 = letter == '%';
    result0 = temp0 || temp1;
    temp0 = letter >= 'a';
    temp1 = letter <= 'z';
    temp0 = temp0 && temp1;
    result0 = temp0 || result0;
    temp0 = letter >= 'A';
    temp1 = letter <= 'Z';
    temp0 = temp0 && temp1;
    check = temp0 || result0;

    if (!check) goto letter_check_done;
    // interpret the code (after the first %)

    // ((lenCodeSoFar == 1) && (codeSoFar[0] == '%'));
    temp0 = lenCodeSoFar == 1;
    temp1 = codeSoFar[0] == '%';
    check = temp0 && temp1;

    if (check) goto just_percentage;
    check = lenCodeSoFar > 0;
    if (check) goto more_formatting_code;
    // else
    goto no_formatting_code;


just_percentage:
    // Just a % sign
    retString[lenRetString] = '%';
    lenRetString++;
    goto interpret_letter_done;

more_formatting_code:
    temp0 = lenCodeSoFar - 1;
    formattingCodeLetter = codeSoFar[temp0];

    /*
     * Find the minimum width (number before the formatting letter)
     */

    minWidth = 0;
    charcounter = 0;
loop6:
    check = charcounter >= (lenCodeSoFar - 1);
    if (check) goto out6;

    // Substract 2 to get the right N for 10^N when converting the
    // number from a string to an int.


    temp0 = lenCodeSoFar - 2;
    tenth = temp0 - charcounter;

    // Check if it's a digit between '0' and '9'

    temp2 = codeSoFar[charcounter]; // NB: temp2 is used below
    temp0 = temp2 >= '0';
    temp1 = temp2 <= '9';
    check = temp0 && temp1;

    if (check) goto is_a_digit;
    // else
    goto not_a_digit;

is_a_digit:

    // Convert digit to int
    digit = (int)(temp2 - '0');

    // digit * 10^tenth

    expo = 0;
loop7:
    check = expo >= tenth;
    if (check) goto out7;
    temp0 = digit * 10;
    digit = temp0;
    expo++;
    goto loop7;
out7:
    minWidth += digit;
    goto digit_check_done;
not_a_digit:
    putchar('E');
    putchar('R');
    putchar('R');
    putchar('O');
    putchar('R');
    putchar(':');
    putchar(' ');
    putchar(codeSoFar[charcounter]);
    putchar('\n');
digit_check_done:

    charcounter++;
    goto loop6;
out6:

    // Check what the formatting letter is
    check = formattingCodeLetter == 'c';
    if (check) goto letter_is_c;
    check = formattingCodeLetter == 'd';
    if (check) goto letter_is_d;
    check = formattingCodeLetter == 's';
    if (check) goto letter_is_s;
    check = formattingCodeLetter == 'x';
    if (check) goto letter_is_x;
    goto letter_is_unknown;

letter_is_c:
    // Get the next arguments
    // must use int instead of char because of how va_arg works in C
    chararg = (char)va_arg(args, int);

    // Add minWidth - 1 spaces in front of the char
    i = 0;
loop8:
    temp0 = minWidth - 1;
    check = (i >= temp0);
    if (check) goto out8;
    retString[lenRetString] = ' ';
    lenRetString++;
    i++;
    goto loop8;
out8:
    // Add char argument to the retString
    retString[lenRetString] = chararg;
    lenRetString++;

    goto letter_done;
letter_is_d:
    // Get the next arguments
    // must use int instead of char because of how va_arg works in C
    intarg = va_arg(args, int);

    /*
     * Build up a string based on the decimal number in intarg
     */

    myitoa(intarg, numberstring, 10);

    // Find the length of numberstring
    i = 0;
loop9:
    check = numberstring[i] == '\0';
    if (check) goto out9;
    i++;
    goto loop9;
out9:

    strlen = i;

    // Add (minWidth-strlen) spaces in front of the number string
    i = 0;
loop10:
    temp0 = minWidth - strlen;
    check = i >= temp0;
    if (check) goto out10;
    retString[lenRetString] = ' ';
    lenRetString++;
    i++;
    goto loop10;
out10:

    // Add number string to the retString
    i = 0;
loop11:
    check = i >= strlen;
    if (check) goto out11;
    retString[lenRetString] = numberstring[i];
    lenRetString++;
    i++;
    goto loop11;
out11:

    goto letter_done;
letter_is_s:

    // Get the next arguments

    // must use int instead of char because of how va_arg works in C
    strarg = (char*)va_arg(args, char*);

    strlen=0;
loop12:
    check = strarg[strlen] == 0;
    if (check) goto out12;
    strlen++;
    goto loop12;
out12:

    // Add (minWidth-strlen) spaces in front of the string
    i = 0;
loop13:
    temp0 = minWidth - strlen;
    check = i >= temp0;
    if (check) goto out13;
    retString[lenRetString] = ' ';
    lenRetString++;
    i++;
    goto loop13;
out13:

    // Add string argument to the retString
    i = 0;
loop14:
    check = i >= strlen;
    if (check) goto out14;
    retString[lenRetString] = strarg[i];
    lenRetString++;
    i++;
    goto loop14;
out14:

    goto letter_done;
letter_is_x:
    // Get the next arguments
    // must use int instead of char because of how va_arg works in C
    intarg = va_arg(args, int);

    /*
     * Build up a string based on the decimal number in intarg
     */

    myitoa(intarg, numberstring, 16);

    // Find the length of numberstring
    i = 0;
loop15:
    check = numberstring[i] == '\0';
    if (check) goto out15;
    i++;
    goto loop15;
out15:

    strlen = i;

    // Add (minWidth-strlen) spaces in front of the number string
    i = 0;
loop16:
    temp0 = minWidth - strlen;
    check = i >= temp0;
    if (check) goto out16;
    retString[lenRetString] = ' ';
    lenRetString++;
    i++;
    goto loop16;
out16:

    // Add number string to the retString
    i = 0;
loop17:
    check = i >= strlen;
    if (check) goto out17;
    retString[lenRetString] = numberstring[i];
    lenRetString++;
    i++;
    goto loop17;
out17:

    goto letter_done;
letter_is_unknown:
    // WEIRDO CODE: &codeSoFar
    putchar('E');
    putchar('R');
    putchar('R');
    putchar('O');
    putchar('R');
    putchar('\n');
letter_done:
    // nop
    //putchar('\n');
    goto interpret_letter_done;
no_formatting_code:
    // "WOOT, NO FORMATTING"
    putchar('E');
    putchar('R');
    putchar('R');
    putchar('O');
    putchar('R');
    putchar('\n');
interpret_letter_done:
    // Done interpreting a formatting string
    inCode = -1; // false
    goto letter_check_done;

just_a_regular_letter:
    check = letter == '%';
    if (check) goto letter_is_percentage;
    // else
    goto letter_is_not_percentage;

letter_is_percentage:
    // Entering code
    inCode = 0; // true
    lenCodeSoFar = 0;
    goto letter_check_done;

letter_is_not_percentage:
    check = letter == 0;
    if (check) goto letter_check_done;

    retString[lenRetString] = letter;
    lenRetString++;
    retString[lenRetString] = 0; // This must be here!

letter_check_done:

    // Stop at end of string
    check = letter == 0;
    if (check) goto out4;

    // Next letter of res
    lettercursor++;
    letter = format[lettercursor];

    goto loop4;
out4:

    va_end(args);

    // Copy the result into res before returning

    i = 0;
loop5:
    res[i] = retString[i];
    check = i == lenRetString;
    if (check) goto out5;
    i++;
    goto loop5;
out5:

    res[lenRetString] = 0;

    return lenRetString;
}
