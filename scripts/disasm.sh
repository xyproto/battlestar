#!/bin/sh
objdump -d -s --disassembler-options=intel "$@"
