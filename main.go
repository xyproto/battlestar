package main

import (
	"strings"
	"strconv"
	"io/ioutil"
	"os"
	"log"
	"fmt"
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

func (tok Token) String() string {
	if tok.t == REGISTER {
		return "register" + tok.value
	} else if tok.t == ASSIGNMENT {
		return "assignment[" + tok.value + "]"
	} else if tok.t == VALUE {
		return "value[" + tok.value + "]"
	} else if tok.t == VALID_NAME {
		return "name[" + tok.value + "]"
	} else if tok.t == SEP {
		return ";"
	} else if tok.t == UNKNOWN {
		return "?[" + tok.value + "]"
	} else if tok.t == KEYWORD {
		return "keyword[" + tok.value + "]"
	} else if tok.t == STRING {
		return "string[" + tok.value + "]"
	}
	log.Fatalln("Error when serializing: Unfamiliar token: " + tok.value + " (?)")
	return "!?"
}

type Statement []Token

var registers = []string{"eax", "ebx", "ecx", "esp", "ebp"}
var operators = []string{"="}
var keywords = []string{"fun", "int", "ret", "const"}

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
	combined := letters + upper + digits + special
	if !(strings.Contains(letters, string(s[0])) || strings.Contains(upper, string(s[0]))) {
		// Does not start with a letter
		return false
	}
	for _, letter := range s {
		// If not a letter, digit or valid special character, it's not a valid name
		if !(strings.Contains(combined, string(letter))) {
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
				collected += word
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

func (st Statement) String() string {
	if len(st) == 0 {
		log.Fatalln("Error: Empty statement:", st)
		return ""
	} else if (st[0].t == KEYWORD) && (st[0].value == "int") { // interrrupt call
		asmcode := "--- call interrupt 0x" + st[1].value + " ---\n"
		// Check the number of parameters
		if len(st) > 6 {
			log.Fatalln("Error: Too many parameters for interrupt call:\n", st)
		}
		// Store each of the parameters to the appropriate registers
		var (
			reg string
			n string
			comment string
		)
		for i := 2; i < len(st); i++ {
			reg = []string{"eax", "ebx", "ecx", "edx"}[i-2]
			n = strconv.Itoa(i-2)
			if (i-2) == 0 {
				comment = "function call: " + st[i].value
			} else {
				if st[i].t == VALUE {
					comment = "parameter #" + n + " is " + st[i].value
				} else {
					comment = "parameter #" + n + " is " + "&" + st[i].value
				}
			}
			asmcode += "mov " + reg + ", " + st[i].value + "\t; " + comment + "\n"
		}
		// Add the interrupt call
		if (st[1].t == VALUE) {
			asmcode += "int 0x" + st[1].value + "\t; perform the call"
			return asmcode
		}
		log.Fatalln("Error: Need a (hexadecimal) interrupt number to call:\n", st)
	} else if (st[0].t == KEYWORD) && (st[0].value == "const") && (len(st) == 4) { // constant data
		if (st[1].t == VALID_NAME) && (st[2].t == ASSIGNMENT) && (st[3].t == STRING) {
			// For the .DATA section (recognized by the keyword)
			return st[1].value + ":\tdb " + st[3].value + "\t; constant string"
		}
		log.Fatalln("Error: Invalid parameters for constant string statement:\n", st)
	} else if (st[0].t == KEYWORD) && (st[0].value == "const") && (len(st) == 5) { // constant data with len
		// For the .DATA section (recognized by the keyword)
		return ";TODO\t; constant: length of " + st[4].value
		log.Fatalln("Error: Invalid parameters for constant string len statement:\n", st)
	} else if (st[0].t == KEYWORD) && (st[0].value == "ret") {
		fmt.Println(st[1])
		return ";TODO: RET"
	} else if len(st) == 3 {
		if (st[0].t == REGISTER) && (st[1].t == ASSIGNMENT) && (st[2].t == VALUE) {
			return "mov " + st[0].value + ", " + st[2].value + "\t; " + st[0].value + " " + st[1].value + " " + st[2].value
		} else {
			log.Fatalln("Error: Uknown type of statement, but familiar layout:\n", st)
		}
	} else if len(st) == 2 {
		if (st[0].t == KEYWORD) && (st[1].t == VALID_NAME) {
			if st[0].value == "fun" {
				asmcode := "--- function " + st[1].value + " ---\n"
				asmcode += "global go.main." + st[1].value + "\t; make label available to linker (Go)\n"
				asmcode += "go.main." + st[1].value + ":\n"
				asmcode += "\n; -- setup stack frame\n"
				asmcode += "push rbp\t; save old base pointer\n"
				asmcode += "mov rbp, rsp\t; use stack pointer as new base pointer\n"
				return asmcode
			} else {
				log.Fatalln("Error: Unknown keyword:", st[0].value)
			}
		}
	}
	log.Println("Error: Unfamiliar statement layout: ")
	for _, token := range []Token(st) {
		log.Print(token.value + " ")
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
					log.Println("CONSTANT")
					if debug {
						log.Println("CONSTANT", asmline)
					}
					constants += asmline + "\n"
				} else {
					if debug {
						log.Println("STATEMENT", statement)
					}
					if debug2 {
						log.Println("ASM", asmline)
					}
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
	log.Println("Battlestar compiler")
	log.Println("Version 0.1")
	log.Println("Alexander Rødseth, 2014")
	log.Println("MIT licensed")

	// Read code from stdin and output 64-bit assembly code
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err == nil {
		tokens := tokenize(string(bytes), true)
		log.Println("--- Done tokenizing ---")
		constants, asmcode := TokensToAssembly(tokens, true, false)
		fmt.Println("SECTION .data\n")
		fmt.Println(constants + "\n")
		fmt.Println("SECTION .text\n")
		fmt.Println(asmcode)
	}
}
