package battlestarlib

import (
	"log"
	"strconv"
	"strings"
)

var (
	// TODO: Several of the operators (and the difference between signed/unsigned)
	operators = []string{"=", "+=", "-=", "*=", "/=", "&=", "|=", "^=", "->", "<<<", ">>>", "<<", ">>", "<->", "==>", "<=="}

	comparisons = []string{"==", "!=", "<", ">", "<=", ">="}

	// TODO: "use" and make the bootable kernel work somehow
	keywords = []string{"fun", "ret", "const", "call", "extern", "end", "bootable", "counter", "address", "value", "loopwrite", "rawloop", "loop", "break", "continue", "use", "asm", "mem", "readbyte", "readword", "readdouble", "membyte", "memword", "memdouble", "var", "write", "noret"}

	// TODO: "read"
	builtins = []string{"len", "int", "exit", "halt", "chr", "print", "read", "syscall"} // built-in functions

	reserved = []string{"funparam", "sysparam", "a", "b", "c", "d"} // built-in lists that can be accessed with [index], or register aliases
)

func qualifier(s string) bool {
	switch s {
	case "byte", "BYTE", "word", "WORD", "dword", "DWORD", "ptr", "PTR", "short", "SHORT", "long", "LONG":
		return true
	}
	return false
}

func validName(s string) bool {
	if len(s) == 0 {
		return false
	}
	if qualifier(s) {
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

// Replace \n, \t, \r and \0 with the appropriate values
func stringReplacements(s string) string {
	rtable := map[string]int{"\\t": 9, "\\n": 10, "\\r": 13, "\\0": 0}
	for key, value := range rtable {
		if strings.Contains(s, key) {
			if strings.Contains(s, key+"\"") {
				s = strings.Replace(s, key+"\"", "\", "+strconv.Itoa(value), -1)
			} else {
				s = strings.Replace(s, key, "\", "+strconv.Itoa(value)+", \"", -1)
			}
		}
	}
	return s
}

func isValue(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (config *TargetConfig) reservedAndValue(st Statement) string {
	if st[0].Value == "funparam" {
		paramoffset, err := strconv.Atoi(st[1].Value)
		if err != nil {
			log.Fatalln("Error: Invalid offset for", st[0].Value+":", st[1].Value)
		}
		return config.paramnum2reg(paramoffset)
	} else if st[0].Value == "sysparam" {
		paramoffset, err := strconv.Atoi(st[1].Value)
		if err != nil {
			log.Fatalln("Error: Invalid offset for", st[0].Value+":", st[1].Value)
		}
		if paramoffset >= len(config.interruptParameterRegisters) {
			log.Fatalln("Error: Invalid offset for", st[0].Value+":", st[1].Value, "(too high)")
		}
		return config.interruptParameterRegisters[paramoffset]
	} else {
		// TODO: Implement support for other lists
		log.Fatalln("Error: Can only handle \"funparam\" and \"sysparam\" reserved words.")
	}
	log.Fatalln("Error: Unable to handle reserved word and value:", st[0].Value, st[1].Value)
	return ""
}
