Battlestar
==========

[![Build Status](https://travis-ci.org/xyproto/battlestar.svg?branch=master)](https://travis-ci.org/xyproto/battlestar)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/xyproto/battlestar/master/LICENSE)
[![Report Card](https://img.shields.io/badge/go_report-A-brightgreen.svg?style=flat)](http://goreportcard.com/report/xyproto/battlestar)


What is Battlestar?
-------------------

* A work in progress.
* Created for fun and for the educational process.
* A programming language specifically for 64-bit x86 Linux, 32-bit x86 Linux and 16-bit x86 DOS.
* Subset of assembly with an alternative syntax and with support for inline C.
* The indended purpose is for writing 4k and 64k demoscene demos.


Quick start
-----------

Build and install Battlestar, build the samples and run the "life" sample:

* `make; sudo make devinstall; make samples; cd life; ./life.sh`

This requires DosBox, Go, Yasm and GCC.

gcc-go is not supported yet.


Features and limitations
------------------------

* The resulting executables are tiny!
* "hello world" is only *174* bytes for 32-bit Linux (when using sstrip from elfkickers). (238 bytes for 64-bit Linux, 31 bytes for 16-bit DOS)
* It's possible to write an operating system / kernel with only one source file.
* Full support for inline C (by utilizing gcc).
* C and Battlestar code can exist in the same source file and calls can be made both ways.
* Battlestar programs compiles almost instantly.
* Programs can be run like scripts by including this line at the top: ```#!/usr/bin/bts```
* Interrupts can be called with the same syntax for both 32-bit and 64-bit x86 on Linux.
* Also supports 16-bit x86 with DosBox.
* The intermediate assembly is fully commented.
* No register allocator, just an alternative assembly syntax.


Sample program
--------------

This is a 16-bit x86 program, for DOS:

```c
// "Life"
// The original was written by "melezov" (http://256bytes.untergrund.net/demo/334)

fun main
    al = 0x13               // set graphics mode (mode 13h)
                            // 320x200, 256 colors, one byte per pixel
    int 10

    stack -> sp             // pop  sp
    stack -> b              // pop  bx
    stack -> ds             // pop  ds

    ds -> es                // push ds, pop es

    al = 62
    ch = 0xFA
    loopwrite               // rep stosb

    loop

        di <<< 3           // rotate left 3

        di -= 7            // subtraction
        di ^= 2            // xor

        al = readbyte di   // read byte from memory
        al += [di+321]     // add value at [di+321] (pixel on the line below)
        al /= 2

        di -> stack
        write               // stosb
        write

        di += 0x13E
        write
        write

        stack -> di

    end // loops forever

end
```

In progress
-----------

* macOS support
* Reimplementing 16-bit demoscene demos without using any inline assembly
* See TODO


Installation
------------------

Make sure Go, Yasm and GCC are installed.

Installation:

`sudo make PREFIX=/usr install`

For development, install soft links instead:

`sudo make install-dev`

Build all the samples:

`make samples`


Build dependencies
------------------
* go >= 1.3


Runtime dependencies
--------------------
* yasm


Optional runtime dependencies
-----------------------------
* gcc (for inline C support)
* elftools/sstrip (for even smaller binaries)
* binutils (for disassembling with objdump)
* dosbox (for running 16-bit executables) (only GCC 4.9 and up supports compiling to 16-bit with -m16)
* SDL 2 (must be compiled and installed manually if on Red Hat 6)
* tcc (for even smaller binaries, in many cases)

External links
--------------

* Battlestar programs on [Rosetta Code](https://rosettacode.org/wiki/Category:Battlestar#mw-pages)

General info
------------

* Version: 0.6.1
* License: MIT
* Author: Alexander F. Rødseth
