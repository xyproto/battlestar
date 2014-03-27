Project 1 - return an error code
--------------------------------

1. Create a file named first.bts

2. Type in:

   exit 3

3. Compile and run:

   bts run first.bts

4. See that the error code was indeed 3:

   echo $?


Project 2 - hello world
-----------------------

1. Create a file named hello.bts

2. Type in:

   const hi = "Hello, World!", 10
   write(hi)

3. Compile and run:

   bts run hello.bts

4. See:

   "Hello, World!"


Project 3 - compile, run and clean
----------------------------------

1. Compile hello.bts from project 2:

   bts build hello.bts

2. See that is is tiny:

   bts size hello.bts

3. Run it:

   ./hello

4. Clean up by removing all the generated files:

   bts clean


Project 4 - run as a script
---------------------------

1. Add this at the top of hello.bts:

   #!/usr/bin/bts

2. Make it executable:

   chmod +x hello.bts

3. Run it as a script:

   ./hello.bts


