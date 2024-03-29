## TODO

- [ ] Support aarch64
- [ ] Make bottle99 and fibonacci work on macOS + aarch64
- [ ] Fix the issue with defining string constants like this: "..", 0 or like this: 46, 46, 0
- [ ] Fix the issue with mul / imul in the spongy sample, see "make todo".
- [ ] Make it possible to use "->" and "<-" with variables, like for the stack.
- [ ] Support for adc, cwd, jz and jnz (use the loop label automatically)
- [ ] Reimplement more 16-bit demoscene demos.
- [ ] Require yasm or nasm, not just yasm.
- [ ] Need a way to differentiate between 8-bit, 16-bit, 32-bit and 64-bit numbers and parameters.
      Something like: x:8 ? x.8 ? x8 ? For every variable name ending with 8, 16, 32 or 64, the number becomes the bits?
- [ ] Add support for Kolibri OS http://wiki.kolibrios.org/wiki/Writing_applications_for_KolibriOS
- [ ] Manpage
- [ ] String blocks can consists of: constants, immediate strings, registers (interpreted as ASCII numbers) and numbers (interpreted as strings)
- [ ] Add internal state that keeps track of what the values in the various registers are used for.
- [ ] Create a keyword "keep" or "protected". If a register is protected, it will be pushed and popped at every function call.
- [ ] Consider making cx/ecx/rcx protected by default in every loop.
- [ ] Consider removing "rawloop".
- [ ] Make "use" work with C libraries. For including a library+include files. Either automatic inclusion of the right .h files with #include.
      OR Automatic linking of libraries and inclusion of header files in a block of inline C when functions are not found. Like printf, sin, cos and SDL_CreateRenderer.
- [ ] "use sdl2" (for example) at the top will add the right compilation and linking flags, using pkg-config
- [ ] Add nested loops
- [ ] Solve rosetta code tasks + programming language benchmark game tasks.
- [ ] Implement a more expressive sub-language and/or inline Go/Julia/Lua/IO
- [ ] Align comments. Drop all the "\t" use in the code.
- [ ] Add functions for checking which token combinations qualifies for which treatment.
- [ ] Use the token module that comes with Go
- [ ] Write code for matching { and }, so that void main() { is not confused by a premature }
- [ ] Built in Quaternions, Matrices, Lists, Vectors and doubles. No 16-bit float?
- [ ] Local variables (.bss section or on the stack? ) or heap?
        - [ ] .data for constants
    - [ ] .bss for uninitialized variables
    - [ ] the heap for local variables that will not use too much memory?
    - [ ] the stack for the rest?
- [ ] Add and fix local variables. Use the .bss section. (Remember to include .bss when building kernels)
- [ ] Return values (not only eax/rax etc)
- [ ] Remove the [ebp-8] local variable testing code.
- [ ] Add the assembly version of printf.
- [ ] Make it easy to link with OpenGL or SDL2.
- [ ] Parse trees instead of token lists, for function calls.
- [ ] Test on Cygwin on Wine as well, possibly change the uses of uname
- [ ] Add a -debug=true flag that includes the C std library and makes it possible to use printf. Drop the assembly version of printf, not needed. Should also compile with -O1 -g etc.
- [ ] Incorporate more ideas (and gcc flags) from this stackoverflow answer: http://stackoverflow.com/a/10552160/131264
- [ ] Add a way to replace one line of tokens with several lines of tokens.
- [ ] Create a standard library that contains platform-dependent battlestar-functions
- [ ] http://nickdesaulniers.github.io/blog/2014/04/18/lets-write-some-x86-64/

### Maybe

- [ ] Syntax for linking with libraries, that connects with pkg-config? "use sdl" should use the right link and inclusions.


### Other ideas

Import a C library, as found by pkg-config. The symbols will then be available as zlib::SYMBOL_NAME:

import zlib


l = [1, 2, 3, 4, 5] // a list

for x in l do
    blabla
end


Ignore comments when // is not the first thing on the line

Built in quaternions

Built in support for SDL2

Use a subset of the ruby syntax


Inline various languages:

c
  hello("hello from C");
end

python
  hello("hello from python")
end
