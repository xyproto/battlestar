void hi(char* msg, int len); // External battlestar function

void c_hi() {
  char *c = "hi ";
  hi(c, 3);
  hi("you\n", 4);
}
