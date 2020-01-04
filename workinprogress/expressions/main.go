package main

import (
	"fmt"
	"log"
	"strings"
)

type ParseTree struct {
	left  *ParseTree
	right *ParseTree
	value string
}

var (
	order = map[string]int{"^": 2, "*": 1, "/": 1, "+": 0, "-": 0}
)

const (
	MAX_ORDERNUM = 2
)

//
// Somehwat cumbersome way to define functions that check if a string only contains a set of letters
//

type OnlyFunction func(string) bool

// Generate a function that checks if a given string consists only of the letters given here.
func OnlyGenerator(checkletters string) OnlyFunction {
	return func(word string) bool {
		for _, letter := range word {
			if !strings.Contains(checkletters, string(letter)) {
				return false
			}
		}
		return true
	}
}

// Combine two OnlyFunctions into one, where the result is or-ed.
func OrGenerator(a, b OnlyFunction) OnlyFunction {
	return func(word string) bool {
		return a(word) || b(word)
	}
}

func KeysAsString(m map[string]int) string {
	s := ""
	for k := range m {
		s += k
	}
	return s
}

var (
	OnlyLetters         = OnlyGenerator("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	OnlyDigits          = OnlyGenerator("0123456789")
	OnlyOperators       = OnlyGenerator(KeysAsString(order))
	OnlyLettersOrDigits = OrGenerator(OnlyLetters, OnlyDigits)
)

//
// Functions that actually does something
//

// Surround all operators by space, then remove all duplicate spaces
func spaceify(s string) string {
	ns := s
	log.Println("ORIGINAL!", ns)
	for _, operator := range KeysAsString(order) {
		ns = strings.Replace(ns, string(operator), " "+string(operator)+" ", -1)
	}
	log.Println("SPACED!", ns)
	for strings.Contains(ns, "  ") {
		ns = strings.Replace(ns, "  ", " ", -1)
	}
	log.Println("FIXED!", ns)
	return ns
}

// Parse a string as a flat list of tokens
func flatparse(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	var tokens []string
	word := ""
	inparen := 0
	for _, letter := range s {
		switch letter {
		case ' ':
			if inparen == 0 {
				// Register the word
				tokens = append(tokens, word)
				word = ""
				continue
			} else {
				word += string(letter)
			}
		case '(':
			inparen += 1
		case ')':
			inparen -= 1
		default:
			word += string(letter)
		}
	}
	// Register the word
	tokens = append(tokens, word)
	return tokens
}

// Parse a string into a ParseTree
func parse(s string) ParseTree {
	if s == "" {
		return ParseTree{nil, nil, s}
	}
	tokens := flatparse(s)
	if len(tokens) == 2 {
		log.Fatalln("Error: Unexpected number of elements:", tokens)
	}
	if len(tokens) >= 3 {
		// Check if there are two operators or two non-operators in a row
		var (
			operator     bool
			lastoperator bool
		)
		for i, token := range tokens {
			lastoperator = operator
			operator = OnlyOperators(token)
			if i == 0 {
				// Skip the first round
				continue
			}
			if operator && lastoperator {
				log.Fatalln("Error: Two operators in a row:", s)
			}
			if (!operator) && (!lastoperator) {
				log.Fatalln("Error: Two non-operators in a row:", s)
			}
		}
	}

	// Find a sorted list (by operator presedence) of the operators
	// Store in the operators map (token number is the key, the string is the value)
	operators := map[int]string{}
	tokenpos_operators := []int{} // List for keeping the order
	for ordernum := 0; ordernum <= MAX_ORDERNUM; ordernum++ {
		for i, token := range tokens {
			if order[token] == ordernum {
				if OnlyOperators(token) {
					operators[i] = token
					tokenpos_operators = append(tokenpos_operators, i)
				}
			}
		}
	}

	// List the operators in order
	//for _, tokenpos := range tokenpos_operators {
	//	fmt.Println(operators[tokenpos], "has order", order[operators[tokenpos]], "and tokenpos", tokenpos)
	//}

	// If there are no operators, return a branch with no branches below
	if len(tokenpos_operators) == 0 {
		return ParseTree{nil, nil, s}
	}

	// Split at the operator with the highest presedence
	tokenpos := tokenpos_operators[0]
	value := operators[tokenpos]

	left_string := strings.Join(tokens[:tokenpos], " ")
	right_string := strings.Join(tokens[tokenpos+1:], " ")

	log.Println("operator", value)
	log.Println("left string", left_string)
	log.Println("right string", right_string)

	left_parse_tree := parse(left_string)
	right_parse_tree := parse(right_string)

	return ParseTree{&left_parse_tree, &right_parse_tree, value}
}

// TODO: Build with a buffer instead of adding to a string
func buildstring(pt *ParseTree, indent int) string {
	s := pt.value + "\n"
	if pt.left != nil {
		s += strings.Repeat("  ", indent)
		s += buildstring(pt.left, indent+1)
	}
	if pt.right != nil {
		s += strings.Repeat("  ", indent)
		s += buildstring(pt.right, indent+1)
	}
	return s
}

// TODO: Build with a buffer instead of adding to a string
func buildexpression(pt *ParseTree) string {
	var left, right, operator string
	s := ""
	if pt.left != nil {
		left = buildexpression(pt.left)
		if OnlyLettersOrDigits(left) {
			s += left
		} else {
			s += "(" + left + ")"
		}

	}
	operator = pt.value
	if OnlyOperators(operator) {
		s += " " + operator + " "
	} else {
		s += operator
	}
	if pt.right != nil {
		right = buildexpression(pt.right)
		if OnlyLettersOrDigits(right) {
			s += right
		} else {
			s += "(" + right + ")"
		}
	}
	return s
}

func (pt ParseTree) String() string {
	return buildexpression(&pt)
}

func (pt ParseTree) Tree() string {
	return buildstring(&pt, 1)
}

func main() {
	//s := "x * (2 + a)"
	//s := "3 + x * (2 + a)"

	// TODO: Don't ignore incoming parantheses
	s := "x^2 * (2+x^3*77-x^2/123) + (1^2*177*(3+3+x-z))"

	pt := parse(spaceify(s))
	fmt.Println(pt.Tree())
	fmt.Println(s)
	fmt.Println(pt.String())
}
