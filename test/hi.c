void hi(char* msg); // External assembly function

void hi_twice() {
  char *c = "there";
  hi(c);
  hi("you");
}
