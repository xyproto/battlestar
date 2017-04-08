package main

// TODO: Add line numbers to the error messages and make them parseable by editors and IDEs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// Global variables
var (
	// 32-bit (i686), 64-bit (x86_64) or 16-bit (i386)
	platformBits = 32

	// Is this a bootable kernel? (declared with "bootable" at the top)
	bootable_kernel = false

	// OS X or Linux
	osx = false
)

func main() {
	name := "Battlestar"
	version := "0.5"
	log.Println(name + " " + version)

	ps := NewProgramState()

	// TODO: Add an option for not adding an exit function
	// TODO: Automatically discover 32-bit/64-bit and Linux/OS X
	// TODO: ARM support

	// Check for -bits=32 or -bits=64 (default)
	platformBitsArg := flag.Int("bits", 64, "Output 32-bit or 64-bit x86 assembly")
	// Check for -osx=true or -osx=false (default)
	osxArg := flag.Bool("osx", false, "On OS X?")
	// Assembly output file
	asmfileArg := flag.String("o", "", "Assembly output file")
	// C output file
	cfileArg := flag.String("oc", "", "C output file")
	// Input file
	btsfileArg := flag.String("f", "", "BTS source file")
	// Is it not a standalone program, but a component? (just the .o file is needed)
	componentArg := flag.Bool("c", false, "Component, not a standalone program")

	flag.Parse()

	platformBits = *platformBitsArg
	osx = *osxArg
	asmfile := *asmfileArg
	cfile := *cfileArg
	btsfile := *btsfileArg
	component := *componentArg

	if flag.Arg(0) != "" {
		btsfile = flag.Arg(0)
	}

	if btsfile == "" {
		log.Fatalln("Abort: a source filename is needed. Provide one with -f or as the first argument.")
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
		if temptokens := tokenize(string(bytes), " "); (len(temptokens) > 2) && (temptokens[0].t == KEYWORD) && (temptokens[0].value == "bootable") && (temptokens[1].t == SEP) {
			bootable = true
			asmdata += fmt.Sprintf("bits %d\n", platformBits)
		} else {
			// Header for regular programs
			asmdata += fmt.Sprintf("bits %d\n", platformBits)
		}

		// Check if platformBits is valid
		if !hasi([]int{16, 32, 64}, platformBits) {
			log.Fatalln("Error: Unsupported bit size:", platformBits)
		}

		init_interrupt_parameter_registers(platformBits)

		btsCode := addExternMainIfMissing(string(bytes))
		tokens := addExitTokenIfMissing(tokenize(btsCode, " "))
		log.Println("--- Done tokenizing ---")
		constants, asmcode := TokensToAssembly(tokens, true, false, ps)
		if constants != "" {
			asmdata += "section .data\n"
			asmdata += constants + "\n"
		}
		if platformBits == 16 {
			asmdata += "org 0x100\n"
		}
		if !bootable {
			asmdata += "\n"
			asmdata += "section .text\n"
		}
		if platformBits == 16 {
			// If there are defined functions, jump over the definitions and start at
			// the main/_start function. If there is a main function, jump to the
			// linker start function. If not, just start at the top.
			// TODO: This is a quick fix. Don't depend on the comment, find a better way.
			if strings.Count(asmcode, "; name of the function") > 1 {
				if strings.Contains(asmcode, "\nmain:") {
					asmdata += "jmp " + linker_start_function + "\n"
				}
			}
		}
		if asmcode != "" {
			if component {
				asmdata += asmcode + "\n"
			} else {
				asmdata += addStartingPointIfMissing(asmcode, ps) + "\n"
			}
			if bootable {
				reg := "esp"
				if platformBits == 64 {
					reg = "rsp"
				}
				asmdata = strings.Replace(asmdata, "; starting point of the program\n", "; starting point of the program\n\tmov "+reg+", stack_top\t; set the "+reg+" register to the top of the stack (special case for bootable kernels)\n", 1)
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
