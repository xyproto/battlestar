package battlestarlib

import (
	"log"
	"strings"
)

// These are constants that each represent a different type of token
const (
	REGISTER       = 0 // a register, like "rax" or "di"
	ASSIGNMENT     = 1
	VALUE          = 2
	KEYWORD        = 3
	BUILTIN        = 4
	VALIDNAME      = 5
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
	ARROW          = 18
	MEMEXP         = 19 // memory expressions, like [di+321]
	ASMLABEL       = 20
	ROL            = 21 // rotate left instruction
	ROR            = 22 // rotate right instruction
	SEGOFS         = 23 // segment:offset for 16-bit assembly
	CONCAT         = 24
	SHL            = 25 // shift left
	SHR            = 26 // shift right
	QUAL           = 27 // qualifier, like BYTE or WORD
	XCHG           = 28
	OUT            = 29
	IN             = 30
	SEP            = 127 // statement separator
	UNKNOWN        = 255
)

var (
	debug          = true
	tokenDebug     = false
	newTokensDebug = true

	tokenToString = TokenDescriptions{REGISTER: "register", ASSIGNMENT: "assignment", VALUE: "value", VALIDNAME: "name", SEP: ";", UNKNOWN: "?", KEYWORD: "keyword", STRING: "string", BUILTIN: "built-in", DISREGARD: "disregard", RESERVED: "reserved", VARIABLE: "variable", ADDITION: "addition", SUBTRACTION: "subtraction", MULTIPLICATION: "multiplication", DIVISION: "division", COMPARISON: "comparison", ARROW: "stack operation", MEMEXP: "address expression", ASMLABEL: "assembly label", AND: "and", XOR: "xor", OR: "or", ROL: "rol", ROR: "ror", CONCAT: "concatenation", SEGOFS: "segment+offset", SHL: "shl", SHR: "shr", QUAL: "qualifier", XCHG: "xchg", OUT: "out", IN: "in"}
	// see also the top of language.go, when adding tokens
)

type (
	// TokenType is one of the above token constants
	TokenType int

	// Token contains everything needed to know about a parsed token
	Token struct {
		T     TokenType
		Value string
		Line  uint
		extra string // Used when coverting from register to string
	}

	// TokenDescriptions is a map from int to string
	TokenDescriptions map[TokenType]string

	// Statement is a slice of tokens
	Statement []Token
)

// Check if a given map has a given key
func haskey(sm map[TokenType]string, key TokenType) bool {
	_, present := sm[key]
	return present
}

// Represent a Token as a string
func (tok Token) String() string {
	if tok.T == SEP {
		return ";"
	} else if haskey(tokenToString, tok.T) {
		return tokenToString[tok.T] + ":" + tok.Value
	}
	log.Fatalln("Error: Unfamiliar token when representing as string: " + tok.Value)
	return "!?"
}

// Represent a TokenType as a string
func (toktyp TokenType) String() string {
	if toktyp == SEP {
		return ";"
	} else if haskey(tokenToString, toktyp) {
		return tokenToString[toktyp]
	}
	log.Fatalln("Error when serializing: Unfamiliar token type when representing tokentype as string: ", int(toktyp))
	return "!?"
}

// Split a string into more tokens and tokenize them
func (config *TargetConfig) retokenize(word string, sep string) []Token {
	var newtokens []Token
	words := strings.Split(word, sep)
	for _, s := range words {
		tokens := config.Tokenize(s, sep)
		//log.Println("RETOKEN", tokens)
		for _, t := range tokens {
			if t.T != SEP {
				newtokens = append(newtokens, t)
			}
		}
	}
	return newtokens
}

func logtoken(tok Token) {
	if tokenDebug {
		log.Println("TOKEN", tok)
	}
}

func lognewtokens(tokens []Token) {
	if newTokensDebug {
		log.Println("NEWTOKENS", tokens)
	}
}

// Tokenize a string
func (config *TargetConfig) Tokenize(program, sep string) []Token {
	statements := maps(maps(strings.Split(program, "\n"), strings.TrimSpace), removecomments)
	tokens := make([]Token, 0)
	var (
		t           Token
		instring    = false // Have we encountered a " for any given statement?
		constexpr   = false // Are we in a constant expression?
		varexpr     = false // Are we in a variable expression?
		collected   string  // Collected string, until end of line
		inlineC     = false // Are we in parts of the code that are inline_c ... end ?
		cBlock      = false // Are we in parts of the code that are void ... } ?
		statementnr uint
	)
	for statementnrInt, statement := range statements {
		// TODO: Use line number instead of statement number (but statement numbers are better than nothing)
		statementnr = uint(statementnrInt)
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
			cBlock = true
			// Skip the start of this type of inline C, don't include "void" as a token
			continue
		} else if inlineC && (words[0] == "end") {
			if debug {
				log.Println("Found the end of inline C block")
			}
			// End both types of blocks when "end" is encountered
			inlineC = false
			cBlock = false
			// Skip the end keyword of this type of inline C block, don't include "end" as a token
			continue
		} else if cBlock && (words[0] == "}") {
			if debug {
				log.Println("Found the } of void C block")
			}
			cBlock = false
			// Skip the } keyword of this type of inline C block, don't include "}" as a token
			continue
		} else if words[0] == "inline_c" {
			if debug {
				log.Println("Found inline_c, starting inline C block")
			}
			inlineC = true
			// Skip the start of this type of inline C, don't include "inline_c" as a token
			continue
		} else if inlineC || cBlock {
			// In a block of inline code, skip and don't include as tokens
			// log.Println("Skipping when tokenizing:", words)
			continue
		}
		// If we are defining a constant, ease up on tokenizing the rest of the line recursively
		if words[0] == "const" {
			constexpr = true
		} else if words[0] == "var" {
			varexpr = true
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
				logtoken(t)
			} else if has(comparisons, word) {
				t = Token{COMPARISON, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
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
				case "==>":
					tokentype = OUT
				case "<==":
					tokentype = IN
				case "<<<":
					tokentype = ROL
				case ">>>":
					tokentype = ROR
				case "<<":
					tokentype = SHL
				case ">>":
					tokentype = SHR
				case "->":
					tokentype = ARROW
				case "<->":
					tokentype = XCHG
				default:
					log.Fatalln("Error: Unhandled operator:", word)
				}
				t = Token{tokentype, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if has(keywords, word) {
				t = Token{KEYWORD, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if has(builtins, word) {
				t = Token{BUILTIN, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if has(reserved, word) {
				if has([]string{"a", "b", "c", "d"}, word) {
					reg := word
					switch config.PlatformBits {
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
				logtoken(t)
			} else if isValue(word) {
				t = Token{VALUE, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if word == "_" {
				t = Token{DISREGARD, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if strings.HasSuffix(word, "++") {
				firstpart := word[:len(word)-2]
				newtokens := config.retokenize(firstpart+" += 1", " ")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
			} else if strings.HasSuffix(word, "--") {
				firstpart := word[:len(word)-2]
				newtokens := config.retokenize(firstpart+" -= 1", " ")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
			} else if validName(word) {
				t = Token{VALIDNAME, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if qualifier(word) {
				t = Token{QUAL, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if strings.Contains(word, "(") {
				newtokens := config.retokenize(word, "(")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
			} else if strings.Contains(word, ")") {
				newtokens := config.retokenize(word, ")")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
			} else if strings.Contains(word, "[") {
				newtokens := config.retokenize(word, "[")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
			} else if strings.Contains(word, "]") {
				newtokens := config.retokenize(word, "]")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
			} else if (!constexpr && !varexpr) && strings.Contains(word, ",") {
				newtokens := config.retokenize(word, ",")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
			} else if strings.Contains(word, "..") {
				newtokens := config.retokenize(word, "..")
				tokens = append(tokens, newtokens...)
				lognewtokens(newtokens)
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
				logtoken(t)
			} else if strings.Contains(word, "+") {
				// Assume it's an address, like bp+5
				t = Token{MEMEXP, "[" + word + "]", statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if strings.Contains(word, "-") {
				// Assume it's an address, like si-0x6
				t = Token{MEMEXP, "[" + word + "]", statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if strings.HasSuffix(word, ":") {
				t = Token{ASMLABEL, word, statementnr, ""}
				tokens = append(tokens, t)
				logtoken(t)
			} else if strings.Count(word, ":") == 1 {
				regs := strings.Split(word, ":")
				if has(registers, regs[0]) && has(registers, regs[1]) {
					// segment:offset
					t = Token{SEGOFS, "[" + word + "]", statementnr, ""}
					tokens = append(tokens, t)
					logtoken(t)
				} else {
					log.Fatalln("Unrecognized segment:offset token:", word)
				}
			} else {
				log.Println("TOKEN", word, "unknown")
				log.Fatalln("Error: Unrecognized token:", word)
				return tokens
			}
		}
		if instring {
			if debug {
				log.Println("EXITING STRING AT END OF STATEMENT")
				log.Println("STRING:", collected)
			}
			t = Token{STRING, stringReplacements(collected), statementnr, ""}
			tokens = append(tokens, t)
			instring = false
			collected = ""
		}
		t = Token{SEP, ";", statementnr, ""}
		tokens = append(tokens, t)
		constexpr = false
		varexpr = false
	}
	return tokens
}

// Replace built-in function calls with more basic code
// Note that only replacements that can be done within one statement will work!
func (config *TargetConfig) reduce(st Statement, debug bool, ps *ProgramState) Statement {
	for i := 0; i < (len(st) - 1); i++ {
		if (st[i].T == BUILTIN) && (st[i].Value == "len") {
			// The built-in len() function

			var name string
			var tokenType TokenType

			if st[i+1].T == VALIDNAME {
				// len followed by a valid name
				// replace with the length of the given value

				name = st[i+1].Value

				if !has(ps.definedNames, name) {
					log.Fatalln("Error:", name, "is unfamiliar. Can not find length.")
				}

				// TODO: Create a built-in cap() function too
				//if length, ok := ps.variables[name]; ok {
				//	token_type = st[i+1].t

				//	// remove the element at i+1
				//	st = st[:i+1+copy(st[i+1:], st[i+2:])]

				//	// replace len(name) with the capacity
				//	st[i] = Token{token_type, strconv.Itoa(length), st[0].line, ""}
				//} else {

				tokenType = st[i+1].T

				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]

				// replace len(name) with _length_of_name, or [_length_of_name] if it's in .bss
				if _, ok := ps.variables[name]; ok {
					st[i] = Token{tokenType, "[_length_of_" + name + "]", st[0].Line, ""}
				} else {
					st[i] = Token{tokenType, "_length_of_" + name, st[0].Line, ""}
				}
			} else if st[i+1].T == REGISTER {
				var length string
				switch config.PlatformBits {
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
				st[i] = Token{VALUE, length, st[0].Line, ""}
			}

			if debug {
				log.Println("SUCCESSFUL REPLACEMENT WITH", st[i])
			}
		} else if (st[i].T == BUILTIN) && (st[i].Value == "print") && (st[i+1].T == STRING) {
			log.Fatalln("Error: print can only print const strings, not immediate strings")
		} else if (st[i].T == BUILTIN) && (st[i].Value == "print") && ((st[i+1].T == VALIDNAME) || (st[i+1].T == REGISTER)) {
			// replace print(msg) with
			// int(0x80, 4, 1, msg, len(msg)) on 32-bit
			// syscall(1, msg, len(msg)) on 64-bit

			// TODO: Find a way to output additional statements when a statement can be broken into several statements
			//       (Like printing all the arguments to print, one by one)
			//for _, token := range st[1:] {
			//    Use token instead of st[i+1]
			//}

			var (
				cmd      string
				tokens   []Token
				tokenpos int
				extra    = st[i+1].extra
			)
			switch config.PlatformBits {
			case 64:
				// Special case when printing single bytes, typically from chr(...)
				if st[i+1].Value == "rsp" {
					cmd = "syscall(1, 1, " + st[i+1].Value + ", 1)"
				} else {
					cmd = "syscall(1, 1, " + st[i+1].Value + ", len(" + st[i+1].Value + "))"
				}
				tokens = config.Tokenize(cmd, " ")
				// Position of the token that is to be written
				tokenpos = 3
			case 32:
				// Special case when printing single bytes, typically from chr(...)
				if st[i+1].Value == "esp" {
					cmd = "int(0x80, 4, 1, " + st[i+1].Value + ", 1)"
				} else {
					cmd = "int(0x80, 4, 1, " + st[i+1].Value + ", len(" + st[i+1].Value + "))"
				}
				tokens = config.Tokenize(cmd, " ")
				// Position of the token that is to be written
				tokenpos = 4
			case 16:
				// No simple reduction for 16-bit assembly, it needs several lines of assembly code
				return st
			}

			tokens[tokenpos].extra = extra
			// Replace the current statement with the newly generated tokens
			st = tokens
		} else if (st[i].T == BUILTIN) && (st[i].Value == "chr") && (st[i+1].T == VALIDNAME) {
			log.Fatalln("Error: str of a defined name is to be implemented")
		} else if (st[i].T == BUILTIN) && (st[i].Value == "chr") && (st[i+1].T == REGISTER) {
			register := st[i+1].Value

			// Replace str(register) with a token VALID_NAME with esp/rsp + register name as the value.
			// This is not perfect, but allows us to output register values with a system call.
			switch config.PlatformBits {
			case 64:
				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]
				// replace with the register that contains the address of the string
				st[i] = Token{REGISTER, "rsp", st[0].Line, register} // only a single byte
			case 32:
				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]
				// replace with the register that contains the address of the string
				st[i] = Token{REGISTER, "esp", st[0].Line, register} // only a single byte
			case 16:
				log.Fatalln("Error: chr() is not implemented for 16-bit platforms")
			}
		}
	}
	return st
}

// TokensToAssembly outputs assembly code given a compilation target config and a slice of tokens
func (config *TargetConfig) TokensToAssembly(tokens []Token, debug bool, debug2 bool, ps *ProgramState) (string, string) {
	statement := []Token{}
	asmcode := ""
	constants := ""
	bsscode := ""
	for _, token := range tokens {
		if token.T == SEP {
			if len(statement) > 0 {
				asmline := Statement(statement).String(ps, config)
				if (statement[0].T == KEYWORD) && (statement[0].Value == "const") {
					if strings.Contains(asmline, ":") {
						if debug {
							log.Printf("CONSTANT: \"%s\"\n", strings.Split(asmline, ":")[0])
						}
					} else {
						log.Fatalln("Error: Unfamiliar constant:", asmline)
					}
					constants += asmline + "\n"
				} else if (statement[0].T == KEYWORD) && (statement[0].Value == "var") {
					// Variables are gathered for the .bss section
					bsscode += asmline + "\n"
				} else {
					asmcode += asmline + "\n"
				}
			}
			statement = []Token{}
		} else {
			statement = append(statement, token)
		}
	}
	// Add .bss section, if any
	if bsscode != "" {
		asmcode += "\nsection .bss\n" + bsscode
	}
	return strings.TrimSpace(constants), asmcode
}

// TokenFilter is a function that can check if
// a given token is one of the allowed types
type TokenFilter (func(Token) bool)

// Take a list of permitted token types.
// Return a TokenFilter function.
func only(tokentypes []TokenType) TokenFilter {
	return func(t Token) bool {
		for _, tt := range tokentypes {
			if t.T == tt {
				return true
			}
		}
		return false
	}
}

// Only return the tokens that the given filter function
// returns true for.
func filtertokens(tokens []Token, filterfunc TokenFilter) []Token {
	newtokens := make([]Token, 0)
	for _, t := range tokens {
		if filterfunc(t) {
			newtokens = append(newtokens, t)
		}
	}
	return newtokens
}
