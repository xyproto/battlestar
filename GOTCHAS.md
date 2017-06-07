Gotchas
-------

* Built in function calls and keywords may change the registers. Check with the assembly output.
* Blocks of inline C that starts with `void main(` and ends with `}` can not contain `}` in between
* Blocks of inline C that starts with `inline_c` and ends with `end` can not be within battlestar functions. C functions are provided
* The `write` function changes several registers, including the loop counter (`cx`/`ecx`/`rcx`).
* Not all samples works on OS X yet.
