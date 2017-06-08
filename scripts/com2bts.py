#!/usr/bin/env python
# -*- coding: utf-8 -*-

from subprocess import check_output
from sys import argv, exit

def com2asm(comfilename):
    """Reads a com file and returns a map of label addresses and instructions"""
    disasm = str(check_output(["ndisasm", "-b", "16", "-p", "intel", comfilename]), encoding="utf-8")
    asmlines = {}
    for line in disasm.split("\n")[:-1]:
        asmlines[int(line[:9], 16)] = line[28:]
    # fix a bug in ndisasm
    for i, line in asmlines.items():
        if "fmul to st1" in line:
            asmlines[i] = line.replace("fmul to st1", "fmul st1, st0")
        elif "fmul to st2" in line:
            asmlines[i] = line.replace("fmul to st2", "fmul st2, st0")
        elif "fmul to, st2" in line:
            asmlines[i] = line.replace("fmul to, st2", "fmul st2, st0")
    return asmlines

def shorten(reg):
    if reg in ["ax", "bx", "cx", "dx"]:
        return reg[:1]
    return reg

def shorten_value(value):
    if value.startswith("byte "):
        return value[5:]
    return value

def transform(line, labelmap, labelcounter):
    """Takes an assembly instruction (like "mov ax, 1") and map of labels (addr -> label).
    Returns the corresponding Battlestar code and a newly generated map of labels (addr -> label).
    Also takes and returns a labelcounter."""
    if line.startswith("mov [") and ("-" not in line) and ("+" not in line):
        a, b = line[5:].split(",", 1)
        if "]" in a:
            a = a.split("]")[0]
        if b in ["al", "ah", "bl", "bh", "cl", "ch", "dl", "dh"]:
            return "membyte " + shorten(a) + " = " + shorten(b), {}, labelcounter
        return "memword " + shorten(a) + " = " + shorten(b), {}, labelcounter
    if "[" in line:
        # TODO, also handle asm lines with [ and ]
        return "(asm 16) " + line, {}, labelcounter
    if line.startswith("int "):
        return line.replace("int ", "int(") + ")", {}, labelcounter
    if line.startswith("push word "):
        return shorten(line[10:]) + " -> stack", {}, labelcounter
    if line.startswith("pop word "):
        return "stack -> " + shorten(line[9:]), {}, labelcounter
    if line.startswith("push "):
        return shorten(line[5:]) + " -> stack", {}, labelcounter
    for loopword in ["loop", "jo", "jno", "js", "jns", "je", "jz", "jne", "jnz", "jb", "jnae", "jc", "jnb", "jae", "jnc", "jbe", "jna", "ja", "jnbe", "jl", "jnge", "jge", "jnl", "jle", "jng", "jg", "jnle", "jp", "jpe", "jnp", "jpo", "jcxz", "jecxz", "jmp short"]:
        if line.startswith(loopword + " "):
            adr = int(line.split()[-1], 16)
            if adr in labelmap:
                label = labelmap[adr]
            else:
                label = "bob" + str(labelcounter)
                labelcounter += 1
            return "(asm 16) " + loopword + " " + label, {adr: label}, labelcounter
    if line.startswith("pop "):
        return "stack -> " + shorten(line[4:]), {}, labelcounter
    if line.startswith("mul ") and line.count(" ") == 1:
        return "a *= " + shorten(line[4:]), {}, labelcounter
    if line.startswith("inc "):
        return shorten(line[4:]) + "++", {}, labelcounter
    if line.startswith("dec "):
        return shorten(line[4:]) + "--", {}, labelcounter
    rtable = {"mov": "=", "add": "+=", "sub": "-=", "xor": "^=", "shr": ">>", "shl": "<<"}
    for word, op in rtable.items():
        if line.startswith(word + " "):
            a, b = line[4:].split(",", 1)
            return shorten(a) + " " + op + " " + shorten_value(b), {}, labelcounter
    if line == "stosb":
        return "write", {}, labelcounter
    return "(asm 16) " + line, {}, labelcounter

def pad(s, n):
    """Pad a string with spaces on both sides until long enough"""
    while len(s) < n:
        if len(s) < n:
            s = s + " "
        if len(s) < n:
            s = " " + s
    return s

def com2bts(comfilename, btsfilename):
    bl = []
    bl.append("fun main")
    labelmap = {}
    labelcounter = 1
    # pass  1, for labels
    for addr, line in com2asm(comfilename).items():
        _, lm, labelcounter = transform(line, {}, labelcounter)
        # update the labelmap dictionary with the values in lm
        if lm:
            labelmap.update(lm)
    labelcounter = 1
    # pass 2, for the code
    for addr, line in com2asm(comfilename).items():
        newline, _, labelcounter = transform(line, labelmap, labelcounter)
        if addr in labelmap:
            label = labelmap[addr] + ":" + "  //==[ " + labelmap[addr] + " ]==|"
            newline = "(asm 16) " + label + "\n" + " " *8 + newline
        else:
            newline = " " * 8 + newline
        for i in range(10):
            h = "0x" + str(i)
            # Simplify hex expressions like 0x0 and 0x2 to just 0 and 2
            if newline.endswith(h):
                newline = newline[:-len(h)] + str(int(h, 16))
        bl.append(newline)
    last_asm_line = bl[-1]
    if last_asm_line.endswith("ret"):
        # Remove the final "ret", since this will be added by "end"
        bl = bl[:-1]
        bl.append("end")
    elif "jmp" in last_asm_line:
        # Don't add a final "ret" if the last assembly line is a jump
        bl.append("noret")
    else:
        bl.append("end")
    bl.append("\n// vim: syntax=c ts=4 sw=4 et:")
    open(btsfilename, "w").write("\n".join(bl))

def main():
    if len(argv) < 2:
        print("Usage: com2bts file.com [file.bts]")
        exit(1)
    if len(argv) < 3 or (len(argv) >= 3 and argv[2] == "-"):
        com2bts(argv[1], "/dev/stdout")
    else:
        com2bts(argv[1], argv[2])

if __name__ == "__main__":
    main()
