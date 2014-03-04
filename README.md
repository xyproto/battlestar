Battlestar
==========

What is Battlestar?
-------------------

* Programming language specifically for 64-bit and 32-bit x86 and Linux.
* Consider it to be assembly with friendlier syntax, with support for inline C, for 32-bit and 64-bit x86.
* The resulting executables are tiny.
* Created for fun and for the educational process.
* The indended purpose is for writing 4k and 64k demoscene demos.

In progress
-----------
* OS X support is in progress.

Features
--------------
* Battlestar programs compiles almost instantly and can also be run like scripts by including this line at the top ```#!/usr/bin/bts```
* The resulting executables are tiny (around 600 bytes for hello world)
* C and Battlestar code can exist in the same source file and calls can be made both ways
* Can call interrupts in a way that is independent of if it's on a 32-bit or 64-bit platform

General information
-------------------
* Version: 0.1
* License: MIT
* Author: Alexander Rødseth

Dependencies
------------
* yasm

Optional dependencies
---------------------
* gcc (for inline C support)

