// Generated with Battlestar 0.3, at 2014-10-29 11:17


// thanks http://arjunsreedharan.org/post/82710718100/kernel-101-lets-write-a-kernel
void kmain(void)
{
    char *msg = "Success!";
    char *vidptr = (char*)0xb8000;  // video mem for text begins here.
    unsigned int i = 0;
    // clear all
    while(i < 80 * 25 * 2) {
        // blank character
        vidptr[i] = ' ';
        i++;
        // attribute-byte: light grey on black screen    
        vidptr[i] = 0x07;         
        i++;
    }
    // write the message to every second byte at 0xb8000
    for (i = 0; msg[i] != '\0'; i++) vidptr[i<<1] = msg[i];
    return;
}

