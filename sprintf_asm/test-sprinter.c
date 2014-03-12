#include <stdarg.h>
#include <stdio.h>
#include <string.h>


extern int sprinter (char *res, char *format, ...);


/* check: Check produced string and return value: */

void check (int n, int ret1, int ret2, char *buf1, char *buf2)
{
  if (ret1==ret2 && strcmp(buf1,buf2)==0) {
    printf("Test %2d OK.\n", n);  return;
  }

  if (strcmp(buf1,buf2) != 0) {
    printf("Test %2d: Teksten er \"%s\" men burde vært \"%s\".\n",
	   n, buf1, buf2);
  }
  if (ret1 != ret2) {
    printf("Test %2d: Returverdi er %d men burde vært %d.\n",
	   n, ret1, ret2);
  }
}


/* The main program: */

int main (void)
{
  char t1[2000], t2[2000];
  int r1, r2;

  r1 = sprinter(t1, "");
  r2 = sprintf(t2, "");
  check(1, r1, r2, t1, t2);

  r1 = sprinter(t1, "En lang tekst uten %%-tegn.");
  r2 = sprintf(t2, "En lang tekst uten %%-tegn.");
  check(2, r1, r2, t1, t2);

  r1 = sprinter(t1, "Ett tegn: '%c'.", 'x');
  r2 = sprintf(t2, "Ett tegn: '%c'.", 'x');
  check(3, r1, r2, t1, t2);

  r1 = sprinter(t1, "To tegn: '%c' og '%c'.", 'x', 'y');
  r2 = sprintf(t2, "To tegn: '%c' og '%c'.", 'x', 'y');
  check(4, r1, r2, t1, t2);

  r1 = sprinter(t1, "Tre tegn: '%c', '%2c' og '%4c'.", 'x', 'y', 'z');
  r2 = sprintf(t2, "Tre tegn: '%c', '%2c' og '%4c'.", 'x', 'y', 'z');
  check(5, r1, r2, t1, t2);

  r1 = sprinter(t1, "Lovlige %s er '%%%%', '%1cc', '%%d', '%%s' og '%%x'.",
       "%-spesifikasjoner", '%');
  r2 = sprintf(t2, "Lovlige %s er '%%%%', '%1cc', '%%d', '%%s' og '%%x'.",
       "%-spesifikasjoner", '%');
  check(6, r1, r2, t1, t2);

  r1 = sprinter(t1, "Tre tekster: '%s', '%s' og '%4s'.", 
       "abc...æøå", "alfa -> omega", "");
  r2 = sprintf(t2, "Tre tekster: '%s', '%s' og '%4s'.", 
       "abc...æøå", "alfa -> omega", "");
  check(7, r1, r2, t1, t2);

  r1 = sprinter(t1, "En økning på %d%% er bedre enn en på %d%%!", 27, 8);
  r2 = sprintf(t2, "En økning på %d%% er bedre enn en på %d%%!", 27, 8);
  check(8, r1, r2, t1, t2);

  r1 = sprinter(t1, "Tallet %d ligger i intervallet %d-%d.", 
              -2230, -10000, -1000);
  r2 = sprintf(t2, "Tallet %d ligger i intervallet %d-%d.", 
              -2230, -10000, -1000);
  check(9, r1, r2, t1, t2);

  r1 = sprinter(t1, "Tallene er %0d, %12d og %209d.", 0, 1000, 1000000000);
  r2 = sprintf(t2, "Tallene er %0d, %12d og %209d.", 0, 1000, 1000000000);
  check(10, r1, r2, t1, t2);

  r1 = sprinter(t1, "Det %2s tallet er %1001d.", 
		"største positive", 2147483647);
  r2 = sprintf(t2, "Det %2s tallet er %1001d.", 
	       "største positive", 2147483647);
  check(11, r1, r2, t1, t2);

  r1 = sprinter(t1, "Det nest %s tallet er %d (-%d).", "største negative", 
	      -2147483647, 1);
  r2 = sprintf(t2, "Det nest %s tallet er %d (-%d).", "største negative", 
	      -2147483647, 1);
  check(12, r1, r2, t1, t2);

  r1 = sprinter(t1, "Det aller %s tallet er %d.", "største negative", 
	      -2147483647-1);
  r2 = sprintf(t2, "Det aller %s tallet er %d.", "største negative", 
	      -2147483647-1);
  check(13, r1, r2, t1, t2);

  r1 = sprinter(t1, "%d = 0x%x", 0, 0);
  r2 = sprintf(t2, "%d = 0x%x", 0, 0);
  check(14, r1, r2, t1, t2);

  r1 = sprinter(t1, "%d = 0x%1x", 1234, 1234);
  r2 = sprintf(t2, "%d = 0x%1x", 1234, 1234);
  check(15, r1, r2, t1, t2);

  r1 = sprinter(t1, "%4d = 0x%8x", -88, -88);
  r2 = sprintf(t2, "%4d = 0x%8x", -88, -88);
  check(16, r1, r2, t1, t2);

  return 0;
}
