Battlestar
==========

* Version: 0.2
* License: MIT
* Author: Alexander Rødseth

What is Battlestar?
-------------------

* A work in progress.
* Created for fun and for the educational process.
* A programming language specifically for 64-bit and 32-bit x86 and Linux.
* Subset of assembly with different syntax and support for inline C.
* The indended purpose is for writing 4k and 64k demoscene demos.

Features
--------

* The resulting executables are tiny!
* "hello world" is only *129* bytes on 32-bit Linux (when using sstrip from elfkickers).
* It's possible to write an operating system with only one source file.
* Full support for inline C (by utilizing gcc).
* C and Battlestar code can exist in the same source file and calls can be made both ways.
* Battlestar programs compiles almost instantly.
* Programs can be run like scripts by including this line at the top: ```#!/usr/bin/bts```
* Interrupts can be called with the same syntax for both 32-bit and 64-bit x86 on Linux.

In progress
-----------
* OS X support

Build dependencies
------------------
* go

Runtime dependencies
--------------------
* yasm

Optional runtime dependencies
-----------------------------
* gcc (for inline C support)
* elftools/sstrip (for even smaller binaries)
* binutils (for disassembling with objdump)
* dosbox (for running 16-bit executables)
