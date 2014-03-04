Battlestar
==========

What is Battlestar?
-------------------

* A work in progress.
* Created for fun and for the educational process.
* A programming language specifically for 64-bit and 32-bit x86 and Linux.
* Consider it to be assembly with friendlier syntax, with support for inline C.
* The resulting executables are tiny.
* The indended purpose is for writing 4k and 64k demoscene demos.

In progress
-----------
* OS X support

Features
--------------
* Battlestar programs compiles almost instantly and can also be run like scripts by including this line at the top:
  ```#!/usr/bin/bts```

* The resulting executables are tiny (around 600 bytes for hello world).
* C and Battlestar code can exist in the same source file and calls can be made both ways.
* Interrupts can be called with the same syntax for both 32-bit and 64-bit x86 on Linux.

Build dependencies
------------------
* go

Runtime dependencies
--------------------
* yasm

Optional runtime dependencies
-----------------------------
* gcc (for inline C support)

General information
-------------------
* Version: 0.1
* License: MIT
* Author: Alexander Rødseth

