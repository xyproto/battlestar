package battlestarlib

import (
	"log"
	"strings"
)

// ExtractInlineC retrieves the C code between:
//   inline_c...end
// or
//   void...}
func ExtractInlineC(code string, debug bool) string {
	var (
		clines       string
		inBlockType1 bool
		inBlockType2 bool
		whitespace   = -1 // Where to strip whitespace
	)
	for _, line := range strings.Split(code, "\n") {
		firstword := strings.TrimSpace(removecomments(line))
		if pos := strings.Index(firstword, " "); pos != -1 {
			firstword = firstword[:pos]
		}
		//log.Println("firstword: "+ firstword)
		if !inBlockType2 && !inBlockType1 && (firstword == "inline_c") {
			log.Println("found", firstword, "starting inline_c block")
			inBlockType1 = true
			// Don't include "inline_c" in the inline C code
			continue
		} else if !inBlockType1 && !inBlockType2 && (firstword == "void") {
			log.Println("found", firstword, "starting inBlockType2 block")
			inBlockType2 = true
			// Include "void" in the inline C code
		} else if !inBlockType2 && inBlockType1 && (firstword == "end") {
			log.Println("found", firstword, "ending inline_c block")
			inBlockType1 = false
			// Don't include "end" in the inline C code
			continue
		} else if !inBlockType1 && inBlockType2 && (firstword == "}") {
			log.Println("found", firstword, "ending inBlockType2 block")
			inBlockType2 = false
			// Include "}" in the inline C code
		}

		if !inBlockType1 && !inBlockType2 && (firstword != "}") {
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

// AddExternMainIfMissing will add "extern main" at the top if
// there is a line starting with "void main", or "int main" but no line starting with "extern main".
func (config *TargetConfig) AddExternMainIfMissing(btsCode string) string {
	foundMain := false
	foundExtern := false
	trimline := ""
	for _, line := range strings.Split(btsCode, "\n") {
		trimline = strings.TrimSpace(line)
		if strings.HasPrefix(trimline, "void main") {
			foundMain = true
		} else if strings.HasPrefix(trimline, "int main") {
			foundMain = true
		} else if strings.HasPrefix(trimline, "extern main") {
			foundExtern = true
		}
		if foundMain && foundExtern {
			break
		}
	}
	if foundMain && !foundExtern {
		return "extern main\n" + btsCode
	}
	return btsCode
}

// AddStartingPointIfMissing will check if the resulting code contains a starting point or not,
// and add one if it is missing.
func (config *TargetConfig) AddStartingPointIfMissing(asmcode string, ps *ProgramState) string {
	if strings.Contains(asmcode, "extern "+config.LinkerStartFunction) {
		log.Println("External starting point for linker, not adding one.")
		return asmcode
	}
	if !strings.Contains(asmcode, config.LinkerStartFunction) {
		log.Printf("No %s has been defined, creating one\n", config.LinkerStartFunction)
		var addstring string
		if config.platformBits != 16 {
			addstring += "global " + config.LinkerStartFunction + "\t\t\t; make label available to the linker\n"
		}
		addstring += config.LinkerStartFunction + ":\t\t\t\t; starting point of the program\n"
		if strings.Contains(asmcode, "extern main") {
			//log.Println("External main function, adding starting point that calls it.")
			linenr := uint(strings.Count(asmcode+addstring, "\n") + 5)
			// TODO: Check that this is the correct linenr
			exitStatement := Statement{Token{BUILTIN, "exit", linenr, ""}}
			return asmcode + "\n" + addstring + "\n\tcall main\t\t; call the external main function\n\n" + exitStatement.String(ps, config)
		} else if strings.Contains(asmcode, "\nmain:") {
			//log.Println("...but main has been defined, using that as starting point.")
			// Add "_start:"/"start" right after "main:"
			return strings.Replace(asmcode, "\nmain:", "\n"+addstring+"main:", 1)
		}
		return addstring + "\n" + asmcode

	}
	return asmcode
}

// AddExitTokenIfMissing will check if the code has an exit or ret and
// will add an exit call if it's missing.
func (config *TargetConfig) AddExitTokenIfMissing(tokens []Token) []Token {
	var (
		twolast        []Token
		lasttoken      Token
		filteredTokens = filtertokens(tokens, only([]TokenType{KEYWORD, BUILTIN, VALUE}))
	)
	if len(filteredTokens) >= 2 {
		twolast = filteredTokens[len(filteredTokens)-2:]
		if twolast[1].T == VALUE {
			lasttoken = twolast[0]
		} else {
			lasttoken = twolast[1]
		}
	} else if len(filteredTokens) == 1 {
		lasttoken = filteredTokens[0]
	} else {
		// less than one token, don't add anything
		return tokens
	}

	// If the last keyword token is ret, exit, jmp or end, all is well, return the same tokens
	if (lasttoken.T == KEYWORD) && ((lasttoken.Value == "ret") || (lasttoken.Value == "end") || (lasttoken.Value == "noret")) {
		return tokens
	}

	// If the last builtin token is exit or halt, all is well, return the same tokens
	if (lasttoken.T == BUILTIN) && ((lasttoken.Value == "exit") || (lasttoken.Value == "halt")) {
		return tokens
	}

	// If not, add an exit statement and return
	newtokens := make([]Token, len(tokens)+2)
	copy(newtokens, tokens)

	// TODO: Check that the line nr is correct
	retToken := Token{BUILTIN, "exit", newtokens[len(newtokens)-1].Line, ""}
	newtokens[len(tokens)] = retToken

	// TODO: Check that the line nr is correct
	sepToken := Token{SEP, ";", newtokens[len(newtokens)-1].Line, ""}
	newtokens[len(tokens)+1] = sepToken

	return newtokens
}
