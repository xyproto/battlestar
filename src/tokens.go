package main

import (
	"log"
	"strings"
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
	COMPARISON     = 17
	PUSHPOP        = 18
	COMBINATION    = 19
	ASMLABEL       = 20
	SEP            = 127
	UNKNOWN        = 255
)

var (
	token_to_string = TokenDescriptions{REGISTER: "register", ASSIGNMENT: "assignment", VALUE: "value", VALID_NAME: "name", SEP: ";", UNKNOWN: "?", KEYWORD: "keyword", STRING: "string", BUILTIN: "built-in", DISREGARD: "disregard", RESERVED: "reserved", VARIABLE: "variable", ADDITION: "addition", SUBTRACTION: "subtraction", MULTIPLICATION: "multiplication", DIVISION: "division", COMPARISON: "comparison", PUSHPOP: "stack operation", COMBINATION: "address expression", ASMLABEL: "assembly label"}
)

type (
	TokenType int

	Token struct {
		t     TokenType
		value string
		line  uint
		extra string // Used when coverting from register to string
	}

	TokenDescriptions map[TokenType]string
	Statement         []Token
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
	log.Fatalln("Error: Unfamiliar token when representing as string: " + tok.value)
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
			// In a block of inline code, skip and don't include as tokens
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
				t = Token{REGISTER, word, statementnr, "?"}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if has(comparisons, word) {
				t = Token{COMPARISON, word, statementnr, ""}
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
				case "->":
					tokentype = PUSHPOP
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
				if has([]string{"a", "b", "c", "d"}, word) {
					reg := ""
					switch platform_bits {
					case 64:
						reg = "r" + word
					case 32:
						reg = "e" + word
					}
					reg += "x"
					t = Token{REGISTER, reg, statementnr, ""}
				} else {
					t = Token{RESERVED, word, statementnr, ""}
				}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if is_value(word) {
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
			} else if strings.HasSuffix(word, "++") {
				firstpart := word[:len(word)-2]
				newtokens := retokenize(firstpart+" += 1", " ", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.HasSuffix(word, "--") {
				firstpart := word[:len(word)-2]
				newtokens := retokenize(firstpart+" -= 1", " ", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if is_valid_name(word) {
				t = Token{VALID_NAME, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if strings.Contains(word, "(") {
				newtokens := retokenize(word, "(", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.Contains(word, ")") {
				newtokens := retokenize(word, ")", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.Contains(word, "[") {
				newtokens := retokenize(word, "[", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if strings.Contains(word, "]") {
				newtokens := retokenize(word, "]", debug)
				for _, newtoken := range newtokens {
					tokens = append(tokens, newtoken)
				}
				log.Println("NEWTOKENS", newtokens)
			} else if (!constexpr) && strings.Contains(word, ",") {
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
			} else if strings.Contains(word, "+") {
				// Assume it's an adress, like bp+5
				t = Token{COMBINATION, word, statementnr, ""}
				tokens = append(tokens, t)
				if debug {
					log.Println("TOKEN", t)
				}
			} else if strings.HasSuffix(word, ":") {
				t = Token{ASMLABEL, word, statementnr, ""}
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
			t = Token{STRING, string_replacements(collected), statementnr, ""}
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
// Note that only replacements that can be done within one statement will work!
func reduce(st Statement, debug bool, ps *ProgramState) Statement {
	for i := 0; i < (len(st) - 1); i++ {
		if (st[i].t == BUILTIN) && (st[i].value == "len") {
			// The built-in len() function

			var name string
			var token_type TokenType

			if st[i+1].t == VALID_NAME {
				// len followed by a valid name
				// replace with the length of the given value

				name = st[i+1].value

				if has(ps.variables[ps.in_function], name) {
					// TODO: Find a way to find the length of local variables
					log.Fatalln("Error: finding the length of a local variable is currently not implemented")
				}
				if !has(ps.defined_names, name) {
					log.Fatalln("Error:", name, "is unfamiliar. Can not find length.")
				}

				token_type = st[i+1].t

				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]

				// replace len(name) with _length_of_name
				st[i] = Token{token_type, "_length_of_" + name, st[0].line, ""}
			} else if st[i+1].t == REGISTER {
				var length string
				switch platform_bits {
				case 64:
					length = "4"
				case 32:
					length = "2"
				case 16:
					length = "1"
				}

				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]

				// replace len(register) with the appropriate length
				st[i] = Token{VALUE, length, st[0].line, ""}
			}

			if debug {
				log.Println("SUCCESSFULL REPLACEMENT WITH", st[i])
			}
		} else if (st[i].t == BUILTIN) && (st[i].value == "write") && (st[i+1].t == STRING) {
			log.Fatalln("Error: write can only write const strings, not immediate strings")
		} else if (st[i].t == BUILTIN) && (st[i].value == "write") && ((st[i+1].t == VALID_NAME) || (st[i+1].t == REGISTER)) {
			// replace write(msg) with
			// int(0x80, 4, 1, msg, len(msg)) on 32-bit
			// syscall(1, msg, len(msg)) on 64-bit

			// TODO: Find a way to output additional statements when a statement can be broken into several statements
			//       (Like printing all the arguments to write, one by one)
			//for _, token := range st[1:] {
			//    Use token instead of st[i+1]
			//}

			cmd := ""
			var tokens []Token
			var tokenpos int
			extra := st[i+1].extra
			switch platform_bits {
			case 64:
				cmd = "syscall(1, 1, " + st[i+1].value + ", len(" + st[i+1].value + "))"
				tokens = tokenize(cmd, true, " ")
				// Position of the token that is to be written
				tokenpos = 3
			case 32:
				cmd = "int(0x80, 4, 1, " + st[i+1].value + ", len(" + st[i+1].value + "))"
				tokens = tokenize(cmd, true, " ")
				// Position of the token that is to be written
				tokenpos = 4
			case 16:
				// No simple reduction for 16-bit assembly, it needs several lines of assembly code
				return st
			}

			tokens[tokenpos].extra = extra
			// Replace the current statement with the newly generated tokens
			st = tokens
		} else if (st[i].t == BUILTIN) && (st[i].value == "chr") && (st[i+1].t == VALID_NAME) {
			log.Fatalln("Error: str of a defined name is to be implemented")
		} else if (st[i].t == BUILTIN) && (st[i].value == "chr") && (st[i+1].t == REGISTER) {
			register := st[i+1].value

			// Replace str(register) with a token VALID_NAME with esp/rsp + register name as the value.
			// This is not perfect, but allows us to output register values with a system call.
			switch platform_bits {
			case 64:
				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]
				// replace with the register that contains the address of the string
				st[i] = Token{REGISTER, "rsp", st[0].line, register}
			case 32:
				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]
				// replace with the register that contains the address of the string
				st[i] = Token{REGISTER, "esp", st[0].line, register}
			case 16:
				log.Fatalln("Error: chr() is not implemented for 16-bit platforms")
			}
		}
	}
	return st
}

func TokensToAssembly(tokens []Token, debug bool, debug2 bool, ps *ProgramState) (string, string) {
	statement := []Token{}
	asmcode := ""
	constants := ""
	for _, token := range tokens {
		if token.t == SEP {
			if len(statement) > 0 {
				asmline := Statement(statement).String(ps)
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
