package main

import (
	"strings"
	"strconv"
	"io/ioutil"
	"os"
	"log"
	"fmt"
	"time"
)

type Program string

type TokenType int
const (
	REGISTER = 0
	ASSIGNMENT = 1
	VALUE = 2
	KEYWORD = 3
	VALID_NAME = 4
	STRING = 5
	SEP = 127
	UNKNOWN = 255
)

type Token struct {
	t TokenType
	value string
}

// Global variables
var infunction string // name of the function we are currently in
var defined_names []string // all defined variables/constants/functions

var token_to_string map[string]string = {REGISTER: "register", ASSIGNMENT: "assignment", VALUE: "value", VALID_NAME: "name", SEP: ";", UNKNOWN: "?", KEYWORD: "keyword", STRING: "string"}


// Check if a given map has a given key
func haskey(sm map[string]string, searchkey string) {
	for key, value := range sm {
		if key == string {
			return true
		}
	}
	return false
}

func (tok Token) String() string {
    if tok.t == SEP {
		return ";"
	} else if haskey(token_to_string, tok.t) {
		return token_to_string[tok.t] + "[" + tok.value + "]"
	}
	log.Fatalln("Error when serializing: Unfamiliar token: " + tok.value + " (?)")
	return "!?"
}

type Statement []Token

var registers = []string{"eax", "ebx", "ecx", "edx", "rbp", "rsp", "rax", "rbx", "rcx", "rdx"}
var parameter_registers = []string{"eax", "ebx", "ecx", "edx"}
var operators = []string{"="}
var keywords = []string{"fun", "int", "ret", "const", "len"}

func maps(sl []string, f func (string) string) []string {
	newl := make([]string, len(sl), len(sl))
	for i, element := range sl {
		newl[i] = f(element)
	}
	return newl
}

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
	letters := "abcdefghijklmnopqrstuvwxyz"
	upper := strings.ToUpper(letters)
	digits := "0123456789"
	special := "_"
	combined1 := letters + upper + special
	combined2 := letters + upper + digits + special
	if !(strings.Contains(combined1, string(s[0]))) {
		// Does not start with a letter
		return false
	}
	for _, letter := range s {
		// If not a letter, digit or valid special character, it's not a valid name
		if !(strings.Contains(combined2, string(letter))) {
			return false
		}
	}
	return true
}

func retokenize(word string, sep string, debug bool) []Token {
	var newtokens []Token
	words := strings.Split(word, sep)
	for _, s := range words {
		tokens := tokenize(s, debug)
		//log.Println("RETOKEN", tokens)
		for _, t := range tokens {
			if t.t != SEP {
				newtokens = append(newtokens, t)
			}
		}
	}
	return newtokens
}

func tokenize(program string, debug bool) []Token {
	statements := maps(strings.Split(program, "\n"), strings.TrimSpace)
	tokens := make([]Token, 0, 0)
	var t Token
	var instring bool // Have we encountered a " for any given statement?
	var collected string // Collected string, until end of line
	for _, statement := range statements {
		words := maps(strings.Split(statement, " "), strings.TrimSpace)
		for _, word := range words {
			if word == "" {
				continue
			}
			if instring {
				collected += word + " "
			} else if has(registers, word) {
				if debug {
					log.Println("TOKEN", word, "register")
				}
				t = Token{REGISTER, word}
				tokens = append(tokens, t)
			} else if has(operators, word) {
				if debug {
					log.Println("TOKEN", word, "operator")
				}
				t = Token{ASSIGNMENT, word}
				tokens = append(tokens, t)
			} else if has(keywords, word) {
				if debug {
					log.Println("TOKEN", word, "keyword")
				}
				t = Token{KEYWORD, word}
				tokens = append(tokens, t)
		    } else if i, err := strconv.Atoi(word); err == nil {
				if debug {
					log.Println("TOKEN", i, "value")
				}
				t = Token{VALUE, word}
				tokens = append(tokens, t)
			} else if is_valid_name(word) {
				if debug {
					log.Println("TOKEN", word, "valid name")
				}
				t = Token{VALID_NAME, word}
				tokens = append(tokens, t)
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
			} else if strings.Contains(word, ",") {
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
				collected = word + " "
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
				t = Token{STRING, collected}
				tokens = append(tokens, t)
			}
			instring = false
		}
		t = Token{SEP, ";"}
		tokens = append(tokens, t)
	}
	return tokens
}

func reduce(st Statement) Statement {
	for i := 0; i < (len(st)-1); i++ {
		if (st[i].t == KEYWORD) && (st[i].value == "len") {
			if st[i+1].t == VALID_NAME {
				// len followed by a valid name
				// replace with the length of the given value

				name := st[i+1].value
				token_type := st[i+1].t

				// remove the element at i+1
				st = st[:i+1+copy(st[i+1:], st[i+2:])]

				// replace len(name) with _length_of_name
				st[i] = Token{token_type, "_length_of_" + name}

				log.Println("SUCCESSFULL REPLACEMENT WITH", st[i])
			}
		}
	}
	return st;
}

func (st Statement) String() string {
	reduced := reduce(st)
	if len(reduced) != len(st) {
		return reduced.String()
	}
	if len(st) == 0 {
		log.Fatalln("Error: Empty statement:", st)
		return ""
	} else if (st[0].t == KEYWORD) && (st[0].value == "int") { // interrrupt call
		asmcode := "\t;--- call interrupt 0x" + st[1].value + " ---\n"
		// Check the number of parameters
		if len(st) > 6 {
			log.Println("Error: Too many parameters for interrupt call:")
			for _, t := range st {
				log.Println(t.value)
			}
			os.Exit(1)
		}
		// Store each of the parameters to the appropriate registers
		var (
			reg string
			n string
			comment string
		)
		for i := 2; i < len(st); i++ {
			reg = parameter_registers[i-2]
			n = strconv.Itoa(i-2)
			if (i-2) == 0 {
				comment = "function call: " + st[i].value
			} else {
				if st[i].t == VALUE {
					comment = "parameter #" + n + " is " + st[i].value
				} else {
					if strings.HasPrefix(st[i].value, "_length_of_") {
						comment = "parameter #" + n + " is len(" + st[i].value[11:] + ")"
					} else {
						comment = "parameter #" + n + " is " + "&" + st[i].value
					}
				}
			}
			codeline := "\tmov " + reg + ", " + st[i].value

			// TODO: Find a more elegant way to format the comments in columns
			if len(codeline) > 14 { // for tab formatting
				asmcode += codeline + "\t\t; " + comment + "\n"
			} else {
				asmcode += codeline + "\t\t\t; " + comment + "\n"
			}
		}
		// Add the interrupt call
		if (st[1].t == VALUE) {
			asmcode += "\tint 0x" + st[1].value + "\t\t\t; perform the call\n"
			return asmcode
		}
		log.Fatalln("Error: Need a (hexadecimal) interrupt number to call:\n", st[1].value)
	} else if (st[0].t == KEYWORD) && (st[0].value == "const") && (len(st) == 4) { // constant data
		constname := ""
		if st[1].t == VALID_NAME {
			constname = st[1].value
		} else {
			log.Fatalln(st[1].value, "is not a valid name for a constant")
		}
		asmcode := ""
		if (st[1].t == VALID_NAME) && (st[2].t == ASSIGNMENT) && (st[3].t == STRING) {
			if has(defined_names, constname) {
				log.Fatalln("Error: constant is already defined: " + constname)
			}
			defined_names = append(defined_names, constname)
			// For the .DATA section (recognized by the keyword)
			asmcode += constname + ":\tdb " + st[3].value + "\t\t; constant string\n"
			// Special naming for storing the length for later
			asmcode += "_length_of_" + constname + " equ $ - " + constname + "\t\t; constant string length\n"
			return asmcode
		}
		log.Fatalln("Error: Invalid parameters for constant string statement:\n", st)
	} else if (st[0].t == KEYWORD) && (st[0].value == "ret") {
		asmcode := "\t;--- takedown stack frame ---\n"
		asmcode += "\tmov rsp, rbp\t\t\t; use base pointer as new stack pointer\n"
		asmcode += "\tpop rbp\t\t\t\t; get the old base pointer\n\n"
		if infunction != "" {
			asmcode += "\t;--- return from \"" + infunction + "\" ---\n"
			// Exiting from the function definition
			infunction = ""
		} else {
			asmcode += "\t;--- return ---\n"
		}
		if (len(st) == 2) && (st[1].t == VALUE){
			asmcode += "\tmov rax, " + st[1].value + "\t\t\t; Error code "
			if st[1].value == "0" {
				asmcode += "0 (everything is fine)\n"
			} else {
				asmcode += st[1].value + "\n"
			}
		}
		asmcode += "\tret\t\t\t\t; Return"
		return asmcode
	} else if len(st) == 3 {
		if (st[0].t == REGISTER) && (st[1].t == ASSIGNMENT) && (st[2].t == VALUE) {
			return "mov " + st[0].value + ", " + st[2].value + "\t\t; " + st[0].value + " " + st[1].value + " " + st[2].value
		} else {
			log.Fatalln("Error: Uknown type of statement, but familiar layout:\n", st)
		}
	} else if (len(st) == 2) && (st[0].t == KEYWORD) && (st[1].t == VALID_NAME) && (st[0].value == "fun") {
		asmcode := ";--- function " + st[1].value + " ---\n"
		infunction = st[1].value
		asmcode += "global " + st[1].value + "\t\t\t; make label available to linker (Go)\n"
		asmcode += st[1].value + ":\t\t\t\t; label / name of the function\n\n"
		asmcode += "\t;--- setup stack frame ---\n"
		asmcode += "\tpush rbp\t\t\t; save old base pointer\n"
		asmcode += "\tmov rbp, rsp\t\t\t; use stack pointer as new base pointer\n"
		return asmcode
	} else if (st[0].t == KEYWORD) {
		log.Fatalln("Error: Unknown keyword:", st[0].value)
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
					if debug {
						log.Println("CONSTANT", asmline[:20], "...")
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
	return strings.TrimSpace(constants), strings.TrimSpace(asmcode)
}

func main() {
	name := "Battlestar"
	version := "0.1"
	log.Println(name + " compiler")
	log.Println("Version " + version)
	log.Println("Alexander Rødseth, 2014")
	log.Println("MIT licensed")

	t := time.Now()
	fmt.Printf("; Generated with %s %s, at %s\n\n", name, version, t.String()[:16])

	// TODO: Needed?
	defined_names = make([]string, 0, 0)

	// Read code from stdin and output 64-bit assembly code
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err == nil {
		tokens := tokenize(string(bytes), true)
		log.Println("--- Done tokenizing ---")
		constants, asmcode := TokensToAssembly(tokens, true, false)
		if constants != "" {
			fmt.Println("SECTION .data\n")
			fmt.Println(constants + "\n")
		}
		if asmcode != "" {
			fmt.Println("SECTION .text\n")
			fmt.Println(asmcode + "\n")
		}
	}
}
