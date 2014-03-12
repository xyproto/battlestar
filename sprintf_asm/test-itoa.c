#include <stdarg.h>
#include <stdio.h>
#include <string.h>


extern char* myitoa(int n, char *s, int b);

int main() {
    char* buffer1;
    char buffer2[99] = "abc";

    printf("HI\n");
    buffer1 = myitoa(-32, buffer2, 16);
    printf("Returned buffer: %s\n", buffer1);
    printf("Modified buffer: %s\n", buffer2);
}
