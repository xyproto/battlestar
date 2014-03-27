package main

// TODO: Add line numbers to the error messages and make them parseable by editors and IDEs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	Program   string
	TokenType int
	Token     struct {
		t     TokenType
		value string
		line  uint
		extra string // Used when coverting from register to string
	}
	TokenDescriptions map[TokenType]string
	Statement         []Token
)

const (
	REGISTER       = 0
	ASSIGNMENT     = 1
	VALUE          = 2
	KEYWORD        = 3
	BUILTIN        = 4
	VALID_NAME     = 5
	STRING         = 6
	DISREGARD      = 7
	RESERVED       = 8
	VARIABLE       = 9
	ADDITION       = 10
	SUBTRACTION    = 11
	MULTIPLICATION = 12
	DIVISION       = 13
	AND            = 14
	OR             = 15
	XOR            = 16
	MAGICAL_VALUE  = 17 // Changes depending on the platform
	SEP            = 127
	UNKNOWN        = 255
)

// Global variables
var (
	in_function          string              // name of the function we are currently in
	inline_c             bool                // are we in a block of inline C? (inline_c ... end)
	c_block              bool                // are we in a block of inline C? (void ... })
	defined_names        []string            // all defined variables/constants/functions
	data_not_value_types []string            // all defined constants that are data (x: db 1,2,3,4...)
	variables            map[string][]string // list of variable names per function name
	types                map[string]string   // type of the defined names

	registers = []string{"ah", "al", "bh", "bl", "ch", "cl", "dh", "dl", // 8-bit
		"si", "di", "sp", "bp", "ip", "ax", "bx", "cx", "dx", // 16-bit
		"eax", "ebx", "ecx", "edx", "esi", "edi", "esp", "ebp", "eip", // 32-bit
		"rax", "rbx", "rcx", "rdx", "rsi", "rdi", "rsp", "rbp", "rip", "r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15", "sil", "dil", "spl", "bpl", "xmm8", "xmm9", "xmm10", "xmm11", "xmm12", "xmm13", "xmm14", "xmm15"} // 64-bit

	operators = []string{"=", "+=", "-=", "*=", "/=", "&=", "|="}
	keywords  = []string{"fun", "ret", "const", "call", "extern", "end", "bootable"}
	builtins  = []string{"len", "int", "exit", "halt", "str", "write", "read", "syscall"} // built-in functions
	reserved  = []string{"funparam", "sysparam"}                                             // built-in lists that can be accessed with [index]

	token_to_string = TokenDescriptions{REGISTER: "register", ASSIGNMENT: "assignment", VALUE: "value", VALID_NAME: "name", SEP: ";", UNKNOWN: "?", KEYWORD: "keyword", STRING: "string", BUILTIN: "built-in", DISREGARD: "disregard", RESERVED: "reserved", VARIABLE: "variable", ADDITION: "addition", SUBTRACTION: "subtraction", MULTIPLICATION: "multiplication", DIVISION: "division"}

	// 32-bit (i686) or 64-bit (x86_64)
	platform_bits = 32

	// Is this a bootable kernel? (declared with "bootable" at the top)
	bootable_kernel = false

	// OS X or Linux
	osx = false

	// TODO: Add an option for not adding start symbols
	linker_start_function = "_start"

	// TODO: Add an option for not adding an exit function

	interrupt_parameter_registers []string
)

// Check if a given map has a given key
func haskey(sm map[TokenType]string, key TokenType) bool {
	_, present := sm[key]
	return present
}

// Represent a Token as a string
func (tok Token) String() string {
	if tok.t == SEP {
		return ";"
	} else if haskey(token_to_string, tok.t) {
		return token_to_string[tok.t] + ":" + tok.value
	}
	//log.Fatalln("Error when serializing: Unfamiliar token type when representing token as string: " + tok.value)
	//log.Fatalln("Error: What is this? " + tok.value)
	log.Fatalln("Error: Unfamiliar token: " + tok.value)
	return "!?"
}

// Represent a TokenType as a string
func (toktyp TokenType) String() string {
	if toktyp == SEP {
		return ";"
	} else if haskey(token_to_string, toktyp) {
		return token_to_string[toktyp]
	}
	log.Fatalln("Error when serializing: Unfamiliar token type when representing tokentype as string: ", int(toktyp))
	return "!?"
}

// Maps the function f over a slice of strings
func maps(sl []string, f func(string) string) []string {
	newl := make([]string, len(sl), len(sl))
	for i, element := range sl {
		newl[i] = f(element)
	}
	return newl
}

// Checks if a slice of strings has the given string
func has(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}

func is_valid_name(s string) bool {
	if len(s) == 0 {
		return false
	}
	// TODO: These could be global constants instead
	letters := "abcdefghijklmnopqrstuvwxyz"
	upper := strings.ToUpper(letters)
	digits := "0123456789"
	special := "_·"
	combined := letters + upper + digits + special

	// Does not start with a number
	if strings.Contains(digits, string(s[0])) {
		return false
	}
	// Check that the rest are valid characters
	for _, letter := range s {
		// If not a letter, digit or valid special character, it's not a valid name
		if !(strings.Contains(combined, string(letter))) {
			return false
		}
	}
	// Valid
	return true
}

// Split a string into more tokens and tokenize them
func retokenize(word string, sep string, debug bool) []Token {
	var newtokens []Token
	words := strings.Split(word, sep)
	for _, s := range words {
		tokens := tokenize(s, debug, sep)
		//log.Println("RETOKEN", tokens)
		for _, t := range tokens {
			if t.t != SEP {
				newtokens = append(newtokens, t)
			}
		}
	}
	return newtokens
}

// Remove one line commants, both // and # are ok
func removecomments(s string) string {
	if strings.HasPrefix(s, "//") || strings.HasPrefix(s, "#") {
		return ""
	} else if pos := strings.Index(s, "//"); pos != -1 {
		// Strip away everything after the first // on the line
		return s[:pos]
	} else if pos := strings.Index(s, "#"); pos != -1 {
		// Strip away everything after the first # on the line
		return s[:pos]
	}
	return s
}

// Tokenize a string
func tokenize(program string, debug bool, sep string) []Token {
	statements := maps(maps(strings.Split(program, "\n"), strings.TrimSpace), removecomments)
	tokens := make([]Token, 0, 0)
	var (
		t           Token
		instring    = false // Have we encountered a " for any given statement?
		constexpr   = false // Are we in a constant expression?
		collected   string  // Collected string, until end of line
		inline_c    = false // Are we in parts of the code that are inline_c ... end ?
		c_block     = false // Are we in parts of the code that are void ... } ?
		statementnr uint
	)
	for statementnr_int, statement := range statements {
		// TODO: Use line number instead of statement number (but statement numbers are better than nothing)
		statementnr = uint(statementnr_int)
		words := maps(strings.Split(statement, " "), strings.TrimSpace)

		if len(words) == 0 {
			continue
		}

		if words[0] == "void" {
			if debug {
				log.Println("Found void, starting C block")
			}
			if (len(words) > 1) && (strings.HasPrefix(words[1], "main(")) {
				log.Println("External main function detected.", words[1])
				// Automatically added
				//log.Println("Remember to add \"extern main\" at the top of the file!")
			}
			c_block = true
			// Skip the start of this type of inline C, don't include "void" as a token
			continue
		} else if inline_c && (words[0] == "end") {
			if debug {
				log.Println("Found the end of inline C block")
			}
			// End both types of blocks when "end" is encountered
			inline_c = false
			c_block = false
			// Skip the end keyword of this type of inline C block, don't include "end" as a token
			continue
		} else if c_block && (words[0] == "}") {
			if debug {
				log.Println("Found the } of void C block")
			}
			c_block = false
			// Skip the } keyword of this type of inline C block, don't include "}" as a token
			continue
		} else if words[0] == "inline_c" {
			if debug {
				log.Println("Found inline_c, starting inline C block")
			}
			inline_c = true
			// Skip the start of this type of inline C, don't include "inline_c" as a token
			continue
		} else if inline_c || c_block {
			// In a block of inline C code, skip and don't include as tokens
			// log.Println("Skipping when tokenizing:", words)
			continue
		}

		// If we are defining a constant, ease up on tokenizing the rest of the line recursively
		if words[0] == "const" {
			constexpr = true
		}

		// Tokenize the words
		for _, word := range words {
			if word == "" {
				continue
			}
			// TODO: refactor out code that repeats the same thing
			if instring {
				collected += word + sep
			} else if has(registers, word) {
				t = Token{REGISTER, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if has(operators, word) {
				var tokentype TokenType
				switch word {
				case "=":
					tokentype = ASSIGNMENT
				case "+=":
					tokentype = ADDITION
				case "-=":
					tokentype = SUBTRACTION
				case "*=":
					tokentype = MULTIPLICATION
				case "/=":
					tokentype = DIVISION
				case "&=":
					tokentype = AND
				case "|=":
					tokentype = OR
				case "^=":
					tokentype = XOR
				default:
					log.Fatalln("Error: Unhandled operator:", word)
				}
				t = Token{tokentype, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if has(keywords, word) {
				t = Token{KEYWORD, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if has(builtins, word) {
				t = Token{BUILTIN, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if has(reserved, word) {
				t = Token{RESERVED, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if _, err := strconv.Atoi(word); err == nil {
				t = Token{VALUE, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if word == "_" {
				t = Token{DISREGARD, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if is_valid_name(word) {
				t = Token{VALID_NAME, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if strings.Contains(word, "(") {
				if debug {
					log.Println("RETOKENIZE BECAUSE OF \"(\"")
				}
				newtokens := retokenize(word, "(", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.Contains(word, ")") {
				if debug {
					log.Println("RETOKENIZE BECAUSE OF \")\"")
				}
				newtokens := retokenize(word, ")", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.Contains(word, "[") {
				if debug {
					log.Println("RETOKENIZE BECAUSE OF \"[\"")
				}
				newtokens := retokenize(word, "[", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.Contains(word, "]") {
				if debug {
					log.Println("RETOKENIZE BECAUSE OF \"]\"")
				}
				newtokens := retokenize(word, "]", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if (!constexpr) && strings.Contains(word, ",") {
				if debug {
					log.Println("RETOKENIZE BECAUSE OF \",\"")
				}
				newtokens := retokenize(word, ",", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.Contains(word, "\"") {
				if debug {
					log.Println("TOKEN", word, "is part of a string")
					log.Println("ENTERING STRING")
				}
				instring = true
				// TODO: This does not work, see test02.asm and test03.asm
				if !strings.HasSuffix(word, sep) {
					if len(collected) == 0 {
						collected += word + sep
					} else {
						collected += word + sep
					}
				} else {
					collected += word + sep + "SEPEND"
				}
			} else if strings.Contains("0123456789$", string(word[0])) {
				// Assume it's a value
				t = Token{VALUE, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}

			} else {
				if debug {
					log.Println("TOKEN", word, "unknown")
				}
				log.Fatalln("Error: Unrecognized token:", word)
				return tokens
			}
		}
		if instring {
			if debug {
				log.Println("EXITING STRING AT END OF STATEMENT")
				log.Println("STRING:", collected)
			}
			t = Token{STRING, collected, statementnr, ""}
			tokens = append(tokens, t)
			instring = false
			collected = ""
		}
		t = Token{SEP, ";", statementnr, ""}
		tokens = append(tokens, t)
		constexpr = false
	}
	return tokens
}

// Replace built-in function calls with more basic code
func reduce(st Statement, debug bool) Statement {
	for i := 0; i < (len(st) - 1); i++ {
		// The built-in len() function
		if (st[i].t == BUILTIN) && (st[i].value == "len") && (st[i+1].t == VALID_NAME) {
			// len followed by a valid name
			// replace with the length of the given value

			name := st[i+1].value

			if !has(defined_names, name) {
				log.Fatalln("Error:", name, "is unfamiliar. Can not find length.")
			}
			if has(variables[in_function], name) {
				// TODO: Find a way to find the length of local variables
				log.Fatalln("Error: finding the length of a local variable is currently not implemented")
			}

			token_type := st[i+1].t

			// remove the element at i+1
			st = st[:i+1+copy(st[i+1:], st[i+2:])]

			// replace len(name) with _length_of_name
			st[i] = Token{token_type, "_length_of_" + name, st[0].line, ""}

			if debug {
				log.Println("SUCCESSFULL REPLACEMENT WITH", st[i])
			}
		} else if (st[i].t == BUILTIN) && (st[i].value == "write") && (st[i+1].t == VALID_NAME) {
			// replace write(msg) with
			// int(0x80, 4, 1, msg, len(msg)) on 32-bit
			// syscall(1, msg, len(msg)) on 64-bit
			// TODO: Convert from string to tokens and use them in place of this token
			cmd := ""
			var tokens []Token
			if platform_bits == 32 {
				cmd = "int(0x80, 4, 1, " + st[i+1].value + ", len(" + st[i+1].value + "))"
				tokens = tokenize(cmd, true, " ")
			} else if platform_bits == 64 {
				cmd = "syscall(1, 1, " + st[i+1].value + ", len(" + st[i+1].value + "))"
				tokens = tokenize(cmd, true, " ")
			}
			// Replace the current statement with the newly generated tokens
			st = tokens
		} else if (st[i].t == BUILTIN) && (st[i].value == "str") && (st[i+1].t == VALID_NAME) {
			log.Fatalln("To implement: str(name)")
		} else if (st[i].t == BUILTIN) && (st[i].value == "str") && (st[i+1].t == REGISTER) {
			register := st[i+1].value

			// remove the element at i+1
			st = st[:i+1+copy(st[i+1:], st[i+2:])]

			// Replace str(register) with a token VALID_NAME with esp/rsp + register name as the value.
			// This is not perfect, but allows us to output register values with a system call.
			if platform_bits == 32 {
				st[i] = Token{VALID_NAME, "esp", st[0].line, register}
			} else {
				st[i] = Token{VALID_NAME, "rsp", st[0].line, register}
			}

			if debug {
				log.Println("SUCCESSFULL REPLACEMENT WITH", st[i], "/", register)
			}

		}
	}
	return st
}

func paramnum2reg(num int) string {
	var offset, reg string
	if platform_bits == 32 {
		offset = strconv.Itoa(8 + num*4)
		reg = "ebp"
	} else if platform_bits == 64 {
		offset = strconv.Itoa(num * 8)
		// ref: page 34 at http://people.freebsd.org/~obrien/amd64-elf-abi.pdf (Figure 3.17)
		switch offset {
		case "0":
			return "rdi"
		case "8":
			return "rsi"
		case "16":
			return "rdx"
		case "24":
			return "rcx"
		case "32":
			return "r8"
		case "40":
			return "r9"
		case "48":
			return "xmm0"
		case "64":
			return "xmm1"
		case "72":
			return "xmm2"
		case "80":
			return "xmm3"
		case "88":
			return "xmm4"
		case "96":
			return "xmm5"
		case "104":
			return "xmm6"
		case "112":
			return "xmm7"
		case "120":
			return "xmm8"
		case "128":
			return "xmm9"
		case "136":
			return "xmm10"
		case "144":
			return "xmm11"
		case "152":
			return "xmm12"
		case "160":
			return "xmm13"
		case "168":
			return "xmm14"
		case "176":
			return "xmm15"
			// TODO: Test if the above offsets and registers are correct
		}
		reg = "rbp"
	}
	return "[" + reg + "+" + offset + "]"
}

func reserved_and_value(st Statement) string {
	if st[0].value == "funparam" {
		paramoffset, err := strconv.Atoi(st[1].value)
		if err != nil {
			log.Fatalln("Error: Invalid offset for", st[0].value+":", st[1].value)
		}
		return paramnum2reg(paramoffset)
	} else if st[0].value == "sysparam" {
		paramoffset, err := strconv.Atoi(st[1].value)
		if err != nil {
			log.Fatalln("Error: Invalid offset for", st[0].value+":", st[1].value)
		}
		if paramoffset >= len(interrupt_parameter_registers) {
			log.Fatalln("Error: Invalid offset for", st[0].value+":", st[1].value, "(too high)")
		}
		return interrupt_parameter_registers[paramoffset]
	} else {
		// TODO: Implement support for other lists
		log.Fatalln("Error: Can only handle \"funparam\" and \"sysparam\" reserved words.")
	}
	log.Fatalln("Error: Unable to handle reserved word and value:", st[0].value, st[1].value)
	return ""
}

func syscall_or_interrupt(st Statement, syscall bool) string {
	if syscall {
		// Remove st[-1]
		i := len(st) - 1
		if st[i].t == SEP {
			log.Println("syscall: ignoring: ", st[i]);
			st = st[:i+copy(st[i:], st[i+1:])]
		}
	} else {
		// Remove st[1], if it's not a value
		i := 1
		if st[i].t != VALUE {
		//	log.Println("REMOVING ", st[i]);
			st = st[:i+copy(st[i:], st[i+1:])]
		}
		// Remove st[-1] if it's a SEP
		i = len(st) - 1
		if st[i].t == SEP {
			log.Println("interrupt call: ignoring: ", st[i]);
			st = st[:i+copy(st[i:], st[i+1:])]
		}
	}

	log.Println("system call:")
	for _, token := range st {
		log.Println(token)
	}
	// Debugging
	//if st[3].value != "o" {
	//	log.Fatalln("break at " + st[3].value)
	//}

	// Store each of the parameters to the appropriate registers
	var reg, n, comment, asmcode, precode, postcode string

	// How many tokens to skip before start reading arguments
	preskip := 2
	if syscall {
		preskip = 1
	}

	from_i := preskip     //inclusive
	to_i := len(st) // exclusive
	step_i := 1
	if osx {
		// arguments are pushed in the opposite order for BSD/OSX (32-bit)
		from_i = len(st) - 1 // inclusive
		to_i = 1             // exclusive
		step_i = -1
	}
	first_i := from_i       // 2 for others, len(st)=1 for OSX/BSD
	last_i := to_i - step_i // 2 for OSX/BSD, len(st)-1 for others
	for i := from_i; i != to_i; i += step_i {
		if (i - preskip) >= len(interrupt_parameter_registers) {
			log.Println("Error: Too many parameters for interrupt call:")
			for _, t := range st {
				log.Println(t.value)
			}
			os.Exit(1)
			break
		}
		reg = interrupt_parameter_registers[i-preskip]
		n = strconv.Itoa(i - preskip)
		if (osx && (i == last_i)) || (!osx && (i == first_i)) {
			comment = "function call: " + st[i].value
		} else {
			if st[i].t == VALUE {
				comment = "parameter #" + n + " is " + st[i].value
			} else if st[i].t == REGISTER {
				log.Fatalln("Error: Can't use a register as a parameter to interrupt calls, since they may be overwritten when preparing for the call.\n" +
					"You can, however, use _ as a parameter to use the value in the corresponding register.")
			} else {
				if strings.HasPrefix(st[i].value, "_length_of_") {
					comment = "parameter #" + n + " is len(" + st[i].value[11:] + ")"
				} else {
					if st[i].value == "_" {
						// When _ is given, use the value already in the corresponding register
						comment = "parameter #" + n + " is supposedly already set"
					} else if has(data_not_value_types, st[i].value) {
						comment = "parameter #" + n + " is " + "&" + st[i].value
					} else {
						comment = "parameter #" + n + " is " + st[i].value
						// Already recognized not to be a register
						if platform_bits == 32 {
							if st[i].value == "esp" {
								// Put the value of the register associated with this token at rbp
								// TODO: Figure out why this doesn't work
								precode += "\tsub esp, 4\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
								precode += "\tmov DWORD [esp], " + st[i].extra + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
								postcode += "\tadd esp, 4\t\t\t; move the stack pointer back\n"
							}
						} else if platform_bits == 64 {
							if st[i].value == "rsp" {
								// Put the value of the register associated with this token at rbp
								// TODO: Figure out why this doesn't work
								precode += "\tsub rsp, 8\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
								precode += "\tmov QWORD [rsp], " + st[i].extra + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
								postcode += "\tadd rsp, 8\t\t\t; move the stack pointer back\n"
							}
						}
					}
				}
			}
		}
		codeline := ""
		// Skip parameters/registers that are already set
		if st[i].value == "_" {
			codeline += "\t\t"
		} else {
			if st[i].value == "0" {
				codeline += "\txor " + reg + ", " + reg
			} else {
				// TODO: Remove special case, implement general local variables
				if st[i].value == "x" {
					if platform_bits == 32 {
						codeline += "\tmov " + reg + ", ebp"
						codeline += "\n\tsub " + reg + ", 8"
					} else {
						codeline += "\tmov " + reg + ", rbp"
						codeline += "\n\tsub " + reg + ", 8"
					}
				} else {
					if osx {
						if i == last_i {
							codeline += "\tmov " + reg + ", " + st[i].value
						} else {
							codeline += "\tpush dword " + st[i].value
						}
					} else {
						codeline += "\tmov " + reg + ", " + st[i].value
					}
				}
			}
		}

		// TODO: Find a more elegant way to format the comments in columns
		if len(codeline) >= 16 { // for tab formatting
			asmcode += codeline + "\t\t; " + comment + "\n"
		} else {
			asmcode += codeline + "\t\t\t; " + comment + "\n"
		}
	}
	if syscall {
		precode = "\t;--- system call ---\n" + precode
	} else {
		comment := "\t;--- call interrupt "
		if !strings.HasPrefix(st[1].value, "0x") {
			// add 0x if missing, assume interrupts will always be called by hex
			comment += "0x"
		}
		comment += st[1].value + " ---\n"
		precode = comment + precode
	}
	// Add the interrupt call
	if syscall || (st[1].t == VALUE) {
		if osx {
			// just the way function calls are made on BSD/OSX
			asmcode += "\tsub esp, 4\t\t\t; BSD system call preparation\n"
		}
		if syscall {
			asmcode += "\tsyscall\t\t\t\t; perform the call\n"
		} else {
			// Add 0x if missing, assume interrupts will always be called by hex
			asmcode += "\tint "
			if !strings.HasPrefix(st[1].value, "0x") {
				log.Println("Note: Adding 0x in front of interrupt", st[1].value)
				asmcode += "0x"
			}
			asmcode += st[1].value + "\t\t\t; perform the call\n"
		}
		if osx {
			pushcount := len(st) - 2
			displacement := strconv.Itoa(pushcount * 4) // 4 bytes per push
			asmcode += "\tadd esp, " + displacement + "\t\t\t; BSD system call cleanup\n"
		}
		return precode + asmcode + postcode
	} else {
		log.Fatalln("Error: Need a (hexadecimal) interrupt number to call:\n", st[1].value)
	}
	return ""
}

func (st Statement) String() string {
	debug := true

	reduced := reduce(st, debug)
	if len(reduced) != len(st) {
		return reduced.String()
	}
	if len(st) == 0 {
		log.Fatalln("Error: Empty statement.")
		return ""
	} else if (st[0].t == BUILTIN) && (st[0].value == "int") { // interrrupt call
		return syscall_or_interrupt(st, false)
	} else if (st[0].t == BUILTIN) && (st[0].value == "syscall") {
		return syscall_or_interrupt(st, true)
	} else if (st[0].t == KEYWORD) && (st[0].value == "const") && (len(st) >= 4) { // constant data
		constname := ""
		if st[1].t == VALID_NAME {
			constname = st[1].value
		} else {
			log.Fatalln(st[1].value, "is not a valid name for a constant")
		}
		asmcode := ""
		if (st[1].t == VALID_NAME) && (st[2].t == ASSIGNMENT) && ((st[3].t == STRING) || (st[3].t == VALUE) || (st[3].t == VALID_NAME)) {
			if has(defined_names, constname) {
				log.Fatalln("Error: Can not declare constant, name is already defined: " + constname)
			}
			if (st[3].t == VALID_NAME) && !has(defined_names, st[3].value) {
				log.Fatalln("Error: Can't assign", st[3].value, "to", st[1].value, "because", st[3].value, "is undefined.")
			}
			// Store the name of the declared constant in defined_names
			defined_names = append(defined_names, constname)
			// For the .DATA section (recognized by the keyword)
			if st[3].t == VALUE {
				if platform_bits == 32 {
					asmcode += constname + ":\tdw "
				} else {
					asmcode += constname + ":\tdq "
				}
			} else {
				asmcode += constname + ":\tdb "
				data_not_value_types = append(data_not_value_types, constname)
			}
			for i := 3; i < len(st); i++ {
				asmcode += st[i].value
				// Add a comma between every element but the last one
				if (i + 1) != len(st) {
					asmcode += ", "
				}
			}
			if st[3].t == STRING {
				asmcode += "\t\t; constant string\n"
			} else {
				asmcode += "\t\t; constant value\n"
			}
			// Special naming for storing the length for later
			asmcode += "_length_of_" + constname + " equ $ - " + constname + "\t; size of constant value\n"
			return asmcode
		}
		log.Println("Error: Invalid parameters for constant string statement:")
		for _, t := range st {
			log.Println(t.value)
		}
		os.Exit(1)
	} else if (len(st) > 2) && (st[0].t == VALID_NAME) && (st[1].t == ASSIGNMENT) {
		log.Println("local variable", st[0].value)
		//for _, t := range st[2:] {
		//    log.Println("new value:", t)
		//}
		// TODO: add proper support for 32-bit, 64-bit and local variable offsets
		//       (-8, -12, -16 etc for 32-bit)
		//       (-8, -16, -24 etc for 64-bit)
		// TODO: add the variable name to the proper global maps and slices
		log.Println("WARNING: Local variables are to be implemented, only one is supported for now")
		// TODO: Remember to sub ebp/rbp
		// TODO: Remove this special case and implement general local variables
		codeline := ""
		if platform_bits == 32 {
			codeline += "\tsub ebp, 8\n"
			codeline += "\tmov DWORD [ebp-8], " + st[2].value + "\t\t\t; " + "local variable x!" + "\n"
		} else {
			codeline += "\tsub rbp, 16\n"
			codeline += "\tmov QWORD [rbp-16], " + st[2].value + "\t\t\t; " + "local variable x!" + "\n"
		}
		return codeline
	} else if (st[0].t == BUILTIN) && (st[0].value == "halt") {
		asmcode := "\t; --- full stop ---\n"
		asmcode += "\tcli\t\t; clear interrupts\n"
		asmcode += ".hang:\n"
		asmcode += "\thlt\n"
		asmcode += "\tjmp .hang\t; loop forever\n\n"
		return asmcode
	} else if (st[0].t == BUILTIN) && (st[0].value == "str") {
		log.Fatalln("Error: This usage of str() is yet to be implemented")
	} else if ((st[0].t == KEYWORD) && (st[0].value == "ret")) || ((st[0].t == BUILTIN) && (st[0].value == "exit")) {
		asmcode := ""
		if st[0].value == "ret" {
			if (in_function == "main") || (in_function == linker_start_function) {
				//log.Println("Not taking down stack frame in the main/_start/start function.")
			} else {
				asmcode += "\t;--- takedown stack frame ---\n"
				if platform_bits == 32 {
					asmcode += "\tmov esp, ebp\t\t\t; use base pointer as new stack pointer\n"
					asmcode += "\tpop ebp\t\t\t\t; get the old base pointer\n\n"
				} else {
					asmcode += "\tmov rsp, rbp\t\t\t; use base pointer as new stack pointer\n"
					asmcode += "\tpop rbp\t\t\t\t; get the old base pointer\n\n"
				}
			}
		}
		if in_function != "" {
			if !bootable_kernel {
				asmcode += "\n\t;--- return from \"" + in_function + "\" ---\n"
			}
		} else if st[0].value == "exit" {
			asmcode += "\t;--- exit program ---\n"
		} else {
			asmcode += "\t;--- return ---\n"
		}
		if (len(st) == 2) && (st[1].t == VALUE) && !osx {
			//if platform_bits == 32 {
			//	if st[1].value == "0" {
			//		asmcode += "\t;NEEDED? xor eax, eax\t\t\t; Error code "
			//	} else {
			//		asmcode += "\t;NEEDED? mov eax, " + st[1].value + "\t\t\t; Error code "
			//	}
			//} else {
			//	if st[1].value == "0" {
			//		asmcode += "\t;NEEDED? xor rdi, rdi\t\t\t; Error code "
			//	} else {
			//		asmcode += "\t;NEEDED? mov rdi, " + st[1].value + "\t\t\t; Error code "
			//	}
			//}
			//if st[1].value == "0" {
			//	asmcode += "0 (ok)\n"
			//} else {
			//	asmcode += st[1].value + "\n"
			//}
		}
		if (st[0].value == "exit") || (in_function == "main") || (in_function == linker_start_function) {
			// Not returning from main/_start/start function, but exiting properly
			exit_code := "0"
			if (len(st) == 2) && ((st[1].t == VALUE) || (st[1].t == REGISTER)) {
				exit_code = st[1].value
			}
			if !bootable_kernel {
				if platform_bits == 32 {
					if osx {
						asmcode += "\tpush dword " + exit_code + "\t\t\t; exit code " + exit_code + "\n"
						asmcode += "\tsub esp, 4\t\t\t; the BSD way, push then subtract before calling\n"
					}
					asmcode += "\tmov eax, 1\t\t\t; function call: 1\n"
					if !osx {
						asmcode += "\t"
						if exit_code == "0" {
							asmcode += "xor ebx, ebx"
						} else {
							asmcode += "mov ebx, " + exit_code
						}
						asmcode += "\t\t\t; exit code " + exit_code + "\n"
					}
					asmcode += "\tint 0x80\t\t\t; exit program\n"
				} else {
					asmcode += "\tmov rax, 60\t\t\t; function call: 60\n\t"
					if exit_code == "0" {
						asmcode += "xor rdi, rdi"
					} else {
						asmcode += "mov rdi, " + exit_code
					}
					asmcode += "\t\t\t; return code " + exit_code + "\n"
					asmcode += "\tsyscall\t\t\t\t; exit program\n"
				}
			} else {
				// For bootable kernels, main does not return. Hang instead.
				log.Println("Warning: Bootable kernels has nowhere to return after the main function. You might want to use the \"halt\" builtin at the end of the main function.")
				//asmcode += Statement{Token{BUILTIN, "halt", st[0].line, ""}}.String()
			}
		} else {
			log.Println("IN FUNCTION", in_function)
			// Do not return eax=0/rax=0 if no return value is explicitly provided, by design
			// This allows the return value from the previous call to be returned instead
			asmcode += "\tret\t\t\t\t; Return\n"
		}
		if in_function != "" {
			// Exiting from the function definition
			in_function = ""
		}
		if inline_c {
			// Exiting from inline C
			inline_c = false
			return "; End of inline C block"
		}
		return asmcode
	} else if (st[0].t == REGISTER) || (st[0].t == DISREGARD) && (len(st) == 3) {
		// Statements like "eax = 3" are handled here
		// TODO: Handle all sorts of equivivalents to assembly statements
		if (st[0].t == REGISTER) && (st[1].t == ASSIGNMENT) && (st[2].t == VALUE || st[2].t == VALID_NAME) {
			if st[2].value == "0" {
				return "\txor " + st[0].value + ", " + st[0].value + "\t\t; " + st[0].value + " " + st[1].value + " " + st[2].value
			} else {
				return "\tmov " + st[0].value + ", " + st[2].value + "\t\t; " + st[0].value + " " + st[1].value + " " + st[2].value
			}
		} else if (st[0].t == VALID_NAME) && (st[1].t == ASSIGNMENT) {
			if has(defined_names, st[0].value) {
				log.Fatalln("Error:", st[0].value, "has already been defined")
			} else {
				log.Fatalln("Error:", st[0].value, "is not recognized as a register (and there is no const qualifier). Can't assign.")
			}
		} else if st[0].t == DISREGARD {
			// TODO: If st[2] is a function, one wishes to call it, then disregard afterwards
			return "\t\t\t\t; Disregarding: " + st[2].value + "\n"
		} else if (st[0].t == REGISTER) && (st[1].t == ASSIGNMENT) && (st[2].t == REGISTER) {
			return "\tmov " + st[0].value + ", " + st[2].value + "\t\t\t; " + st[0].value + " " + st[1].value + " " + st[2].value
		} else if (st[0].t == RESERVED) && (st[1].t == VALUE) {
			return reserved_and_value(st[:2])
		} else if (st[0].t == REGISTER) && (st[1].t == ASSIGNMENT) && (st[2].t == RESERVED) && (st[3].t == VALUE) {
			if st[2].value == "param" {
				paramoffset, err := strconv.Atoi(st[3].value)
				if err != nil {
					log.Fatalln("Error: Invalid list offset for", st[2].value+":", st[3].value)
				}
				param_expression := paramnum2reg(paramoffset)
				if len(param_expression) == 3 {
					param_expression += "\t"
				}
				return "\tmov " + st[0].value + ", " + param_expression + "\t\t; fetch function param #" + st[3].value + "\n"
			} else {
				// TODO: Implement support for other lists
				log.Fatalln("Error: Can only handle \"param\" lists.")
			}
		}
		if (st[1].t == ADDITION) && (st[2].t == VALUE) {
			if st[2].value == "1" {
				return "\tinc " + st[0].value + "\t\t\t; " + st[0].value + "++"
			}
			return "\tadd " + st[0].value + ", " + st[2].value + "\t\t; " + st[0].value + " += " + st[2].value
		} else if (st[1].t == SUBTRACTION) && (st[2].t == VALUE) {
			if st[2].value == "1" {
				return "\tdec " + st[0].value + "\t\t\t; " + st[0].value + "--"
			}
			return "\tsub " + st[0].value + ", " + st[2].value + "\t\t; " + st[0].value + " -= " + st[2].value
		} else if (st[1].t == MULTIPLICATION) && (st[2].t == VALUE) {
			// TODO: Don't use a list, write a function that covers the lot
			shifts := []string{"2", "4", "8", "16", "32", "64", "128"}
			if has(shifts, st[2].value) {
				pos := 0
				for i, v := range shifts {
					if v == st[2].value {
						// Found the appropriate shift value
						pos = i + 1
						break
					}
				}
				return "\tshl " + st[0].value + ", " + strconv.Itoa(pos) + "\t\t; " + st[0].value + " *= " + st[2].value
			} else {
				return "\timul " + st[0].value + ", " + st[2].value + "\t\t; " + st[0].value + " *= " + st[2].value
			}
		} else if (st[1].t == DIVISION) && (st[2].t == VALUE) {
			// TODO: Don't use a list, write a function that covers the lot
			shifts := []string{"2", "4", "8", "16", "32", "64", "128"}
			if has(shifts, st[2].value) {
				pos := 0
				for i, v := range shifts {
					if v == st[2].value {
						// Found the appropriate shift value
						pos = i + 1
						break
					}
				}
				return "\tshr " + st[0].value + ", " + strconv.Itoa(pos) + "\t\t; " + st[0].value + " /= " + st[2].value
			} else {
				asmcode := "\n\t;--- signed division: " + st[0].value + " /= " + st[2].value + " ---\n"
				if platform_bits == 32 {
					// Dividing a 64-bit number in edx:eax by the number in ecx. Clearing out edx and only using 32-bit numbers for now.
					// If the register to be divided is rax, do a quicker division than if it's another register
					if st[0].value == "eax" {
						// save ecx
						asmcode += "\tpush ecx\t\t; save ecx\n"
						// save edx
						asmcode += "\tpush edx\t\t; save edx\n"
						// clear edx
						asmcode += "\txor edx, edx\t\t; edx = 0 (32-bit 0:eax instead of 64-bit edx:eax)\n"
						// ecx = st[2].value
						asmcode += "\tmov ecx, " + st[2].value + "\t\t; divisor, ecx = " + st[2].value + "\n"
						// idiv ecx
						asmcode += "\tidiv ecx\t\t\t; eax = edx:eax / ecx\n"
						// restore edx
						asmcode += "\tpop edx\t\t; restore edx\n"
						// restore ecx
						asmcode += "\tpop ecx\t\t; restore ecx\n"
					} else {
						// TODO: if the given register is a different one than eax, ecx and edx,
						//       just divide directly with that register, like for eax above
						// save eax, we know this is not where we assign the result
						asmcode += "\tpush eax\t\t; save eax\n"
						if st[0].value != "ecx" {
							// save ecx
							asmcode += "\tpush ecx\t\t; save ecx\n"
						}
						if st[0].value != "edx" {
							// save edx
							asmcode += "\tpush edx\t\t; save edx\n"
						}
						// copy number to be divided to eax
						asmcode += "\tmov eax, " + st[0].value + "\t\t; dividend, number to be divided\n"
						// clear edx
						asmcode += "\txor edx, edx\t\t; edx = 0 (32-bit 0:eax instead of 64-bit edx:eax)\n"
						// ecx = st[2].value
						asmcode += "\tmov ecx, " + st[2].value + "\t\t; divisor, ecx = " + st[2].value + "\n"
						// eax = edx:eax / ecx
						asmcode += "\tidiv ecx\t\t\t; eax = edx:eax / ecx\n"
						if st[0].value != "edx" {
							// restore edx
							asmcode += "\tpop edx\t\t; restore edx\n"
						}
						if st[0].value != "ecx" {
							// restore ecx
							asmcode += "\tpop ecx\t\t; restore ecx\n"
						}
						// st[0].value = eax
						asmcode += "\tmov " + st[0].value + ", eax\t\t; " + st[0].value + " = eax\n"
						// restore eax
						asmcode += "\tpop eax\t\t; restore eax\n"
					}
					asmcode += "\n"
					return asmcode
				} else {
					// Dividing a 128-bit number in rdx:rax by the number in rcx. Clearing out rdx and only using 64-bit numbers for now.
					// If the register to be divided is rax, do a quicker division than if it's another register
					if st[0].value == "rax" {
						// save rdx
						asmcode += "\tmov r9, rdx\t\t; save rdx\n"
						// clear rdx
						asmcode += "\txor rdx, rdx\t\t; rdx = 0 (64-bit 0:rax instead of 128-bit rdx:rax)\n"
						// mov r8, st[2].value
						asmcode += "\tmov r8, " + st[2].value + "\t\t; divisor, r8 = " + st[2].value + "\n"
						// idiv rax
						asmcode += "\tidiv r8\t\t\t; rax = rdx:rax / r8\n"
						// restore rdx
						asmcode += "\tmov rdx, r9\t\t; restore rdx\n"
					} else {
						log.Println("Note: r8, r9 and r10 will be changed when dividing: " + st[0].value + " /= " + st[2].value)
						// TODO: if the given register is a different one than rax, rcx and rdx,
						//       just divide directly with that register, like for rax above
						// save rax, we know this is not where we assign the result
						asmcode += "\tmov r9, rax\t\t; save rax\n"
						if st[0].value != "rdx" {
							// save rdx
							asmcode += "\tmov r10, rdx\t\t; save rdx\n"
						}
						// copy number to be divided to rax
						asmcode += "\tmov rax, " + st[0].value + "\t\t; dividend, number to be divided\n"
						// xor rdx, rdx
						asmcode += "\txor rdx, rdx\t\t; rdx = 0 (64-bit 0:rax instead of 128-bit rdx:rax)\n"
						// mov rcx, st[2].value
						asmcode += "\tmov r8, " + st[2].value + "\t\t; divisor, r8 = " + st[2].value + "\n"
						// idiv rax
						asmcode += "\tidiv r8\t\t\t; rax = rdx:rax / r8\n"
						if st[0].value != "rdx" {
							// restore rdx
							asmcode += "\tmov rdx, r10\t\t; restore rdx\n"
						}
						// mov st[0].value, rax
						asmcode += "\tmov " + st[0].value + ", rax\t\t; " + st[0].value + " = rax\n"
						// restore rax
						asmcode += "\tmov rax, r9\t\t; restore rax\n"
					}
					return asmcode
				}
			}
		}
		log.Println("Unfamiliar 3-token expression!")
	} else if (len(st) == 4) && (st[0].t == RESERVED) && (st[1].t == VALUE) && (st[2].t == ASSIGNMENT) && (st[3].t == REGISTER) {
		retval := "\tmov " + reserved_and_value(st[:2]) + ", " + st[3].value + "\t\t\t; "
		retval += fmt.Sprintf("%s[%s] = %s\n", st[0].value, st[1].value, st[3].value)
		return retval

	} else if (len(st) == 5) && (st[0].t == RESERVED) && (st[1].t == VALUE) && (st[2].t == ASSIGNMENT) && (st[3].t == RESERVED) && (st[4].t == VALUE) {
		retval := "\tmov " + reserved_and_value(st[:2]) + ", " + reserved_and_value(st[3:]) + "\t\t\t; "
		retval += fmt.Sprintf("%s[%s] = %s[%s]\n", st[0].value, st[1].value, st[3].value, st[4].value)
		return retval
	} else if (len(st) >= 2) && (st[0].t == KEYWORD) && (st[1].t == VALID_NAME) && (st[0].value == "fun") {
		if in_function != "" {
			log.Fatalf("Error: Missing \"ret\" or \"end\"? Already in a function named %s when declaring function %s.\n", in_function, st[1].value)
		}
		asmcode := ";--- function " + st[1].value + " ---\n"
		in_function = st[1].value
		// Store the name of the declared function in defined_names
		if has(defined_names, in_function) {
			log.Fatalln("Error: Can not declare function, name is already defined:", in_function)
		}
		defined_names = append(defined_names, in_function)
		asmcode += "global " + in_function + "\t\t\t; make label available to the linker\n"
		asmcode += in_function + ":\t\t\t\t; name of the function\n\n"
		if (in_function == "main") || (in_function == linker_start_function) {
			//log.Println("Not setting up stack frame in the main/_start/start function.")
			return asmcode
		}
		asmcode += "\t;--- setup stack frame ---\n"
		if platform_bits == 32 {
			asmcode += "\tpush ebp\t\t\t; save old base pointer\n"
			asmcode += "\tmov ebp, esp\t\t\t; use stack pointer as new base pointer\n"
		} else {
			asmcode += "\tpush rbp\t\t\t; save old base pointer\n"
			asmcode += "\tmov rbp, rsp\t\t\t; use stack pointer as new base pointer\n"
		}
		return asmcode
	} else if (st[0].t == KEYWORD) && (st[0].value == "call") && (len(st) == 2) {
		if st[1].t == VALID_NAME {
			return "\t;--- call the \"" + st[1].value + "\" function ---\n\tcall " + st[1].value + "\n"
		} else {
			log.Fatalln("Error: Calling an invalid name:", st[1].value)
		}
		// TODO: Find a shorter format to describe matching tokens.
		// Something along the lines of: if match(st, [KEYWORD:"extern"], 2)
	} else if (st[0].t == KEYWORD) && (st[0].value == "bootable") && (len(st) == 1) {
		bootable_kernel = true
		// This program is supposed to be bootable
		return `
; Thanks to http://wiki.osdev.org/Bare_Bones_with_NASM

; Declare constants used for creating a multiboot header.
MBALIGN     equ  1<<0                   ; align loaded modules on page boundaries
MEMINFO     equ  1<<1                   ; provide memory map
FLAGS       equ  MBALIGN | MEMINFO      ; this is the Multiboot 'flag' field
MAGIC       equ  0x1BADB002             ; 'magic number' lets bootloader find the header
CHECKSUM    equ -(MAGIC + FLAGS)        ; checksum of above, to prove we are multiboot
 
; Declare a header as in the Multiboot Standard. We put this into a special
; section so we can force the header to be in the start of the final program.
; You don't need to understand all these details as it is just magic values that
; is documented in the multiboot standard. The bootloader will search for this
; magic sequence and recognize us as a multiboot kernel.
section .multiboot
align 4
	dd MAGIC
	dd FLAGS
	dd CHECKSUM
 
; Currently the stack pointer register (esp) points at anything and using it may
; cause massive harm. Instead, we'll provide our own stack. We will allocate
; room for a small temporary stack by creating a symbol at the bottom of it,
; then allocating 16384 bytes for it, and finally creating a symbol at the top.
section .bootstrap_stack
align 4
stack_bottom:
times 16384 db 0
stack_top:

section .text
    `
	} else if (st[0].t == KEYWORD) && (st[0].value == "extern") && (len(st) == 2) {
		if st[1].t == VALID_NAME {
			extname := st[1].value
			// Declare the external name
			if has(defined_names, extname) {
				log.Fatalln("Error: Can not declare external symbol, name is already defined: " + extname)
			}
			// Store the name of the declared constant in defined_names
			defined_names = append(defined_names, extname)
			// Return a comment
			return "extern " + extname + "\t\t\t; external symbol\n"
		} else {
			log.Fatalln("Error: extern with invalid name:", st[1].value)
		}
	} else if (st[0].t == KEYWORD) && (st[0].value == "end") && (len(st) == 1) {
		if inline_c {
			inline_c = false
			return "; end of inline C block\n"
		} else if in_function != "" {
			// Return from the function if "end" is encountered
			ret := Token{KEYWORD, "ret", st[0].line, ""}
			newstatement := Statement{ret}
			return newstatement.String()
		} else {
			log.Fatalln("Error: Not in a function or block of inline C, hard to tell what should be ended with \"end\". Statement nr:", st[0].line)
		}
	} else if (st[0].t == VALID_NAME) && (len(st) == 1) {
		// Just a name, assume it's a function call
		if has(defined_names, st[0].value) {
			call := Token{KEYWORD, "call", st[0].line, ""}
			newstatement := Statement{call, st[0]}
			return newstatement.String()
		} else {
			log.Fatalln("Error: No function named:", st[0].value)
		}
	} else if (len(st) > 2) && (st[0].t == VARIABLE) && (st[1].t == ASSIGNMENT) {
		// negative base pointer offset for local variables
		paramoffset := len(variables[in_function]) - 1
		negative_offset := strconv.Itoa(paramoffset*4 + 8)
		reg := ""
		asmcode := ""
		if platform_bits == 32 {
			reg = "ebp"
			asmcode = "\tmov DWORD [" + reg + "-" + negative_offset + "], " + st[2:].String()
		} else {
			negative_offset = strconv.Itoa(paramoffset*8 + 8)
			reg = "rbp"
			asmcode = "\tmov QWORD [" + reg + "-" + negative_offset + "], " + st[2:].String()
		}
		asmcode += "\t\t; local variable #" + strconv.Itoa(paramoffset) + "\n"
		return asmcode
	} else if (st[0].t == KEYWORD) && (st[0].value == "inline_c") {
		inline_c = true
		return "; start of inline C block\n"
	} else if (st[0].t == KEYWORD) && (st[0].value == "const") {
		log.Fatalln("Error: Incomprehensible constant:", st.String())
	} else if st[0].t == BUILTIN {
		log.Fatalln("Error: Unhandled builtin:", st[0].value)
	} else if st[0].t == KEYWORD {
		log.Fatalln("Error: Unhandled keyword:", st[0].value)
	}
	log.Println("Error: Unfamiliar statement layout: ")
	for _, token := range []Token(st) {
		log.Print(token)
	}
	os.Exit(1)
	return ";ERROR"
}

func TokensToAssembly(tokens []Token, debug bool, debug2 bool) (string, string) {
	statement := []Token{}
	asmcode := ""
	constants := ""
	for _, token := range tokens {
		if token.t == SEP {
			if len(statement) > 0 {
				asmline := Statement(statement).String()
				if (statement[0].t == KEYWORD) && (statement[0].value == "const") {
					if strings.Contains(asmline, ":") {
						if debug {
							log.Printf("CONSTANT: \"%s\"\n", strings.Split(asmline, ":")[0])
						}
					} else {
						log.Fatalln("Error: Unfamiliar constant:", asmline)
					}
					constants += asmline + "\n"
				} else {
					asmcode += asmline + "\n"
				}
			}
			statement = []Token{}
		} else {
			statement = append(statement, token)
		}
	}
	return strings.TrimSpace(constants), asmcode
}

func ExtractInlineC(code string, debug bool) string {
	// Fetch the inline C code between the "c" and "end" kewyords
	clines := ""
	inline_c = false
	c_block = false
	whitespace := -1 // Where to strip whitespace
	for _, line := range strings.Split(code, "\n") {
		firstword := strings.TrimSpace(removecomments(line))
		if pos := strings.Index(firstword, " "); pos != -1 {
			firstword = firstword[:pos]
		}
		//log.Println("firstword: "+ firstword)
		if !c_block && !inline_c && (firstword == "inline_c") {
			log.Println("found", firstword, "starting inline_c block")
			inline_c = true
			// Don't include "inline_c" in the inline C code
			continue
		} else if !inline_c && !c_block && (firstword == "void") {
			log.Println("found", firstword, "starting c_block block")
			c_block = true
			// Include "void" in the inline C code
		} else if !c_block && inline_c && (firstword == "end") {
			log.Println("found", firstword, "ending inline_c block")
			inline_c = false
			// Don't include "end" in the inline C code
			continue
		} else if !inline_c && c_block && (firstword == "}") {
			log.Println("found", firstword, "ending c_block block")
			c_block = false
			// Include "}" in the inline C code
		}

		if !inline_c && !c_block && (firstword != "}") {
			// Skip lines that are not in an "inline_c ... end" or "void ... }" block.
			//log.Println("not C, skipping:", line)
			continue
		}

		// Detect whitespace, once and only for some variations
		if whitespace == -1 {
			if strings.HasPrefix(line, "    ") {
				whitespace = 4
			} else if strings.HasPrefix(line, "\t") {
				whitespace = 1
			} else if strings.HasPrefix(line, "  ") {
				whitespace = 2
			} else {
				whitespace = 0
			}
		}
		// Strip whitespace, and check that only whitespace has been stripped
		if (len(line) >= whitespace) && (strings.TrimSpace(line) == strings.TrimSpace(line[whitespace:])) {
			clines += line[whitespace:] + "\n"
		} else {
			clines += line + "\n"
		}
	}
	return clines
}

func add_extern_main_if_missing(bts_code string) string {
	// If there is a line starting with "void main", but no line starting with "extern main",
	// add "extern main" at the top.
	found_main := false
	found_extern := false
	trimline := ""
	for _, line := range strings.Split(bts_code, "\n") {
		trimline = strings.TrimSpace(line)
		if strings.HasPrefix(trimline, "void main") {
			found_main = true
		} else if strings.HasPrefix(trimline, "extern main") {
			found_extern = true
		}
		if found_main && found_extern {
			break
		}
	}
	if found_main && !found_extern {
		return "extern main\n" + bts_code
	}
	return bts_code
}

func add_starting_point_if_missing(asmcode string) string {
	// Check if the resulting code contains a starting point or not
	if strings.Contains(asmcode, "extern "+linker_start_function) {
		log.Println("External starting point for linker, not adding one.")
		return asmcode
	}
	if !strings.Contains(asmcode, linker_start_function) {
		log.Printf("No %s has been defined, creating one\n", linker_start_function)
		addstring := "global " + linker_start_function + "\t\t\t; make label available to the linker\n" + linker_start_function + ":\t\t\t\t; starting point of the program\n"
		if strings.Contains(asmcode, "extern main") {
			//log.Println("External main function, adding starting point that calls it.")
			linenr := uint(strings.Count(asmcode+addstring, "\n") + 5)
			// TODO: Check that this is the correct linenr
			exit_statement := Statement{Token{BUILTIN, "exit", linenr, ""}}
			return asmcode + "\n" + addstring + "\n\tcall main\t\t; call the external main function\n\n" + exit_statement.String()
		} else if strings.Contains(asmcode, "\nmain:") {
			//log.Println("...but main has been defined, using that as starting point.")
			// Add "_start:"/"start" right after "main:"
			return strings.Replace(asmcode, "\nmain:", "\n"+addstring+"main:", 1)
		}
		return addstring + "\n" + asmcode

	}
	return asmcode
}

// Creates and returns a function that can check if
// a given token is one of the allowed types
type TokenFilter (func(Token) bool)

// Take a list of permitted token types.
// Return a TokenFilter function.
func only(tokentypes []TokenType) TokenFilter {
	return func(t Token) bool {
		for _, tt := range tokentypes {
			if t.t == tt {
				return true
			}
		}
		return false
	}
}

// Only return the tokens that the given filter function
// returns true for.
func filtertokens(tokens []Token, filterfunc TokenFilter) []Token {
	newtokens := make([]Token, 0, 0)
	for _, t := range tokens {
		if filterfunc(t) {
			newtokens = append(newtokens, t)
		}
	}
	return newtokens
}

func add_exit_token_if_missing(tokens []Token) []Token {
	var (
		twolast   []Token
		lasttoken Token
	)
	filtered_tokens := filtertokens(tokens, only([]TokenType{KEYWORD, BUILTIN, VALUE}))
	if len(filtered_tokens) >= 2 {
		twolast = filtered_tokens[len(filtered_tokens)-2:]
		if twolast[1].t == VALUE {
			lasttoken = twolast[0]
		} else {
			lasttoken = twolast[1]
		}
	} else if len(filtered_tokens) == 1 {
		lasttoken = filtered_tokens[0]
	} else {
		// less than one token, don't add anything
		return tokens
	}

	// If the last keyword token is ret, exit or end, all is well, return the same tokens
	if (lasttoken.t == KEYWORD) && ((lasttoken.value == "ret") || (lasttoken.value == "end")) {
		return tokens
	}

	// If the last builtin token is exit or halt, all is well, return the same tokens
	if (lasttoken.t == BUILTIN) && ((lasttoken.value == "exit") || (lasttoken.value == "halt")) {
		return tokens
	}

	// If not, add an exit statement and return
	newtokens := make([]Token, len(tokens)+2, len(tokens)+2)
	for i, _ := range tokens {
		newtokens[i] = tokens[i]
	}

	// TODO: Check that the line nr is correct
	ret_token := Token{BUILTIN, "exit", newtokens[len(newtokens)-1].line, ""}
	newtokens[len(tokens)] = ret_token

	// TODO: Check that the line nr is correct
	sep_token := Token{SEP, ";", newtokens[len(newtokens)-1].line, ""}
	newtokens[len(tokens)+1] = sep_token

	return newtokens
}

func main() {
	name := "Battlestar"
	version := "0.1"
	log.Println(name + " compiler")
	log.Println("Version " + version)
	log.Println("Alexander Rødseth")
	log.Println("2014")
	log.Println("MIT licensed")

	// Initialize global maps and slices
	defined_names = make([]string, 0, 0)
	variables = make(map[string][]string)

	// TODO: Automatically discover 32-bit/64-bit and Linux/OS X
	// Check for -bits=32 or -bits=64 (default)
	bits := flag.Int("bits", 64, "Output 32-bit or 64-bit x86 assembly")
	// Check for -osx=true or -osx=false (default)
	is_osx := flag.Bool("osx", false, "On OS X?")
	// Assembly output file
	asm_file := flag.String("o", "", "Assembly output file")
	// C output file
	c_file := flag.String("oc", "", "C output file")
	// Input file
	bts_file := flag.String("f", "", "BTS source file")

	flag.Parse()

	platform_bits = *bits
	osx = *is_osx

	asmfile := *asm_file
	cfile := *c_file
	btsfile := *bts_file

	if flag.Arg(0) != "" {
		btsfile = flag.Arg(0)
	}

	if btsfile == "" {
		log.Fatalln("Abort: An input filename is needed, either by -f or as first argument")
	}

	if asmfile == "" {
		asmfile = btsfile + ".asm"
	}

	if cfile == "" {
		cfile = btsfile + ".c"
	}

	// TODO: Consider adding an option for "start" as well, or a custom
	// start symbol

	if osx {
		linker_start_function = "_main"
	} else {
		linker_start_function = "_start"
	}

	// Assembly file contents
	asmdata := ""

	// C file contents
	cdata := ""

	// Read code from stdin and output 32-bit or 64-bit assembly code
	bytes, err := ioutil.ReadFile(btsfile)
	if err == nil {
		if len(strings.TrimSpace(string(bytes))) == 0 {
			// Empty program
			log.Fatalln("Error: Empty program")
		}

		t := time.Now()
		asmdata += fmt.Sprintf("; Generated with %s %s, at %s\n\n", name, version, t.String()[:16])

		// If "bootable" is the first token
		bootable := false
		if temptokens := tokenize(string(bytes), true, " "); (len(temptokens) > 2) && (temptokens[0].t == KEYWORD) && (temptokens[0].value == "bootable") && (temptokens[1].t == SEP) {
			bootable = true
			// Header for bootable kernels, use 32-bit assembly
			platform_bits = 32
			asmdata += fmt.Sprintf("bits %d\n", platform_bits)
		} else {
			// Header for regular programs
			asmdata += fmt.Sprintf("bits %d\n", platform_bits)
		}

		// Used when calling interrupts (or syscall)
		if platform_bits == 32 {
			interrupt_parameter_registers = []string{"eax", "ebx", "ecx", "edx"}
		} else {
			interrupt_parameter_registers = []string{"rax", "rdi", "rsi", "rdx", "rcx", "r8", "r9"}
		}

		bts_code := add_extern_main_if_missing(string(bytes))
		tokens := add_exit_token_if_missing(tokenize(bts_code, true, " "))
		log.Println("--- Done tokenizing ---")
		constants, asmcode := TokensToAssembly(tokens, true, false)
		if constants != "" {
			asmdata += fmt.Sprintln("section .data") + "\n"
			// TODO: Is Sprintln needed here?
			asmdata += fmt.Sprintln(constants) + "\n"
		}
		if !bootable {
			asmdata += fmt.Sprintln("section .text") + "\n"
		}
		if asmcode != "" {
			asmdata += fmt.Sprintln(add_starting_point_if_missing(asmcode) + "\n")
			if bootable {
				asmdata = strings.Replace(asmdata, "; starting point of the program\n", "; starting point of the program\n\tmov esp, stack_top\t; set the esp register to the top of the stack (special case for bootable kernels)\n", 1)
			}
		}
		ccode := ExtractInlineC(strings.TrimSpace(string(bytes)), true)
		if ccode != "" {
			cdata += fmt.Sprintf("// Generated with %s %s, at %s\n\n", name, version, t.String()[:16])
			cdata += ccode
		}
	}

	log.Println("--- Finalizing ---")

	if asmdata != "" {
		err = ioutil.WriteFile(asmfile, []byte(asmdata), 0644)
		if err != nil {
			log.Fatalln("Error: Unable to write to", asmfile)
		}
		log.Printf("Wrote %s (%d bytes)\n", asmfile, len(asmdata))
	}

	if cdata != "" {
		err = ioutil.WriteFile(cfile, []byte(cdata), 0644)
		if err != nil {
			log.Fatalln("Error: Unable to write to", cfile)
		}
		log.Printf("Wrote %s (%d bytes)\n", cfile, len(cdata))
	}

	log.Println("Done.")
}
