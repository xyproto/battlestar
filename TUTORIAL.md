Project 1 - return an error code
--------------------------------

1. Create a file named `first.bts`

2. Type in:

```c
exit 3
```

3. Compile and run:

```bash
bts run first.bts
```

4. Check that the exit code was indeed `3`:

```bash
echo $?
```


Project 2 - hello world
-----------------------

1. Create a file named `hello.bts`

2. Type in:

```c
const hi = "Hello, World!", 10
write(hi)
```

3. Compile and run:

```bash
bts run hello.bts
```

4. The output should be:

```
"Hello, World!"
```


Project 3 - compile, run and clean
----------------------------------

1. Compile hello.bts from project 2:

```bash
bts build hello.bts
```

2. Observe that the resulting executable is tiny:

```bash
bts size hello.bts
```

3. Run it:

```bash
./hello
```

4. Clean up by removing all the generated files:

```bash
bts clean
```


Project 4 - run as a script
---------------------------

1. Add this at the top of `hello.bts`:

```bash
#!/usr/bin/bts
```

2. Make it executable:

```bash
chmod +x hello.bts
```

3. Run it as a script:

```bash
./hello.bts
```

Project 5 - loop
----------------

1. Create `loop.bts`
2. Type in:

```c
const hi = "Hello there", 10

fun main
  loop 7
    print(hi)
  end
end
```

`"Hello there", 10` is the same as `"Hello there\n"` since `'\n'` has ASCII value 10.
