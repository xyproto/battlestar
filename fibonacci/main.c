#include <stdio.h>

extern int fibb(int n);

int main() {
  printf("fib %d\n", fibb(30));
  return 0;
}

