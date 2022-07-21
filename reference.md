### Quick reference

#### Assign a value to a register

    register = value

example:

    eax = 2

#### Shorthand for using the same register as the platform can offer (16/32/64-bit)

    a = 2

This is translated into "mov eax, 2", for 32-bit x86.

#### Declare a function

    fun name

example:

    fun hello_world

example:

    fun hello
        eax = 2
        ret(0)

#### Return from a function

    ret [retval]

example:

    ret(1)

example:

    ret

#### Declare a constant

    const name = constant

The given value can be a string, comma-separated list of values or a mix of both.

#### Built in functions

    len(name)

Represents the length of a constant, by the given name

#### Comparison starts an if block

    x = 2
    write(message)

#### Operators

The parentheses contains a short explanation and not part of the examples.

Examples:

    a += 2  (addition)
    a -= 2  (subtraction)
    a /= 2  (division - translated to shl/shr when possible)
    a *= 2  (multiplication - translated to shl/shr when possible)
    a |= 2  (or)
    a &= 2  (and)
    a ^= 2  (xor)
    a <<< 2 (rol - rotate bits left)
    a >>> 2 (ror - rotate bits right)

#### Memory access

    a += [di+321]

#### Stack

    ds -> stack     (push ds)
    stack -> es     (pop es)
    ds -> es        (push ds, then pop es)
