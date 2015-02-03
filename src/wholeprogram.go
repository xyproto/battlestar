package main

import (
	"log"
	"strings"
)

var (
	// TODO: Add an option for not adding start symbols
	linker_start_function = "_start"
	inline_c              bool // are we in a block of inline C? (inline_c ... end)
	c_block               bool // are we in a block of inline C? (void ... })
)

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
	// If there is a line starting with "void main", or "int main" but no line starting with "extern main",
	// add "extern main" at the top.
	found_main := false
	found_extern := false
	trimline := ""
	for _, line := range strings.Split(bts_code, "\n") {
		trimline = strings.TrimSpace(line)
		if strings.HasPrefix(trimline, "void main") {
			found_main = true
		} else if strings.HasPrefix(trimline, "int main") {
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

func add_starting_point_if_missing(asmcode string, ps *ProgramState) string {
	// Check if the resulting code contains a starting point or not
	if strings.Contains(asmcode, "extern "+linker_start_function) {
		log.Println("External starting point for linker, not adding one.")
		return asmcode
	}
	if !strings.Contains(asmcode, linker_start_function) {
		log.Printf("No %s has been defined, creating one\n", linker_start_function)
		var addstring string
		if platform_bits != 16 {
			addstring += "global " + linker_start_function + "\t\t\t; make label available to the linker\n"
		}
		addstring += linker_start_function + ":\t\t\t\t; starting point of the program\n"
		if strings.Contains(asmcode, "extern main") {
			//log.Println("External main function, adding starting point that calls it.")
			linenr := uint(strings.Count(asmcode+addstring, "\n") + 5)
			// TODO: Check that this is the correct linenr
			exit_statement := Statement{Token{BUILTIN, "exit", linenr, ""}}
			return asmcode + "\n" + addstring + "\n\tcall main\t\t; call the external main function\n\n" + exit_statement.String(ps)
		} else if strings.Contains(asmcode, "\nmain:") {
			//log.Println("...but main has been defined, using that as starting point.")
			// Add "_start:"/"start" right after "main:"
			return strings.Replace(asmcode, "\nmain:", "\n"+addstring+"main:", 1)
		}
		return addstring + "\n" + asmcode

	}
	return asmcode
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
