package battlestarlib

// TODO Refactor

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	registers = []string{"ah", "al", "bh", "bl", "ch", "cl", "dh", "dl", // 8-bit
		"ax", "bx", "cx", "dx", "si", "di", "sp", "bp", "ip", "cs", "es", "ds", "fs", "gs", "ss", // 16-bit
		"eax", "ebx", "ecx", "edx", "esi", "edi", "esp", "ebp", "eip", // 32-bit
		"rax", "rbx", "rcx", "rdx", "rsi", "rdi", "rsp", "rbp", "rip", "r8", "r9",
		"r10", "r11", "r12", "r13", "r14", "r15", "sil", "dil", "spl", "bpl",
		"xmm8", "xmm9", "xmm10", "xmm11", "xmm12", "xmm13", "xmm14", "xmm15"} // 64-bit

)

// ParseState keeps track of the current state when parsing
type ParseState struct {
	inlineC bool // currently in a block of inline C?
}

// TargetConfig contains information about the current platform and compile target
type TargetConfig struct {
	// PlatformBits should be 16, 32 or 64
	PlatformBits int

	macOS bool

	// BootableKernel should be true if this is not a normal executable but a bootable kernel
	BootableKernel bool

	// LinkerStartFunction is the name of the first function the linker should use, typically "_start"
	LinkerStartFunction string

	// interruptParameterRegisters are the registers that are primarily used when calling interrupts
	interruptParameterRegisters []string
}

// NewTargetConfig returns an new TargetConfig struct with a few options set.
// platformBits should be 16, 32 or 64
// macOS should be true if targeting Darwin / OS X / macOS
// bootableKernel should be set to true if this is for building a bootable kernel
func NewTargetConfig(platformBits int, macOS, bootableKernel bool) (*TargetConfig, error) {
	linkerStartFunction := "_start"
	if macOS {
		linkerStartFunction = "_main"
	}

	// Check if platformBits is valid
	if !hasi([]int{16, 32, 64}, platformBits) {
		return nil, fmt.Errorf("Error: Unsupported bit size: %d", platformBits)
	}

	// Used when calling interrupts (or syscall). Not used for 16-bit platforms.
	var interruptParameterRegisters []string
	if platformBits == 32 {
		interruptParameterRegisters = []string{"eax", "ebx", "ecx", "edx"}
	} else {
		interruptParameterRegisters = []string{"rax", "rdi", "rsi", "rdx", "rcx", "r8", "r9"}
	}

	return &TargetConfig{platformBits, macOS, bootableKernel, linkerStartFunction, interruptParameterRegisters}, nil
}

// is64bit determines if the given register name looks like the 64-bit version of the general purpose registers
func is64bit(reg string) bool {
	// Anything after "rax" (including)
	return pos(registers, reg) >= pos(registers, "rax")
}

// is32bit determines if the given register name looks like the 32-bit version of the general purpose registers
func is32bit(reg string) bool {
	regPos := pos(registers, reg)
	eaxPos := pos(registers, "eax")
	raxPos := pos(registers, "rax")
	// Between "eax" (including) and "rax" (excluding)
	return (eaxPos <= regPos) && (regPos < raxPos)
}

// is16bit determines if the given register name looks like the 16-bit version of the general purpose registers
func is16bit(reg string) bool {
	regPos := pos(registers, reg)
	axPos := pos(registers, "ax")
	eaxPos := pos(registers, "eax")
	// Between "ax" (including) and "eax" (excluding)
	return (axPos <= regPos) && (regPos < eaxPos)
}

// Try to find the 32-bit version of a 64-bit register, or a 16-bit version of a 32-bit register
// If given an empty string, an empty string is returned.
func downgrade(reg string) string {
	if len(reg) == 0 {
		return reg
	}
	if reg[0] == 'r' {
		return "e" + reg[1:]
	}
	if reg[0] == 'e' {
		return reg[1:]
	}
	return reg
}

// Downgrade a register until it is the size of a byte. Requires the string to be non-empty.
func downgradeToByte(reg string) string {
	retval := reg
	if reg[0] == 'r' || reg[0] == 'e' {
		retval = reg[1:]
	}
	return strings.Replace(retval, "x", "l", 1)
}

// Tries to convert a register to a word size register. Requires the string to be non-empty.
func regToWord(reg string) string {
	return upgrade(downgradeToByte(reg))
}

// Tries to convert a register to a double register. Requires the string to be non-empty.
func regToDouble(reg string) string {
	return upgrade(upgrade(downgradeToByte(reg)))
}

func upgrade8bitRegisterTo16bit(reg string) string {
	if len(reg) >= 2 && reg[1] == 'l' {
		return string(reg[0]) + "x"
	} else if len(reg) >= 2 && reg[1] == 'h' {
		return string(reg[0]) + "x"
	}
	return reg
}

// Try to find the 64-bit version of a 32-bit register, or a 32-bit version of a 16-bit register.
// Requires the string to be non-empty.
func upgrade(reg string) string {
	if (reg[0] == 'e') && is64bit("r"+reg[1:]) {
		return "r" + reg[1:]
	}
	if is32bit("e" + reg) {
		return "e" + reg
	}
	if is16bit(upgrade8bitRegisterTo16bit(reg)) {
		return upgrade8bitRegisterTo16bit(reg)
	}
	return reg
}

// Checks if the register is one of the a registers.
func registerA(reg string) bool {
	return (reg == "ax") || (reg == "eax") || (reg == "rax") || (reg == "al") || (reg == "ah")
}

func (config *TargetConfig) paramnum2reg(num int) string {
	var offset, reg string
	switch config.PlatformBits {
	case 64:
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
	case 32:
		offset = strconv.Itoa(8 + num*4)
		reg = "ebp"
	case 16:
		log.Fatalln("Error: PARAMETERS are not implemented for 16-bit assembly, yet")
	}
	return "[" + reg + "+" + offset + "]"
}

func (config *TargetConfig) counterRegister() string {
	switch config.PlatformBits {
	case 16:
		return "cx"
	case 32:
		return "ecx"
	case 64:
		return "rcx"
	default:
		log.Fatalln("Error: Unhandled bit size:", config.PlatformBits)
		return ""
	}
}

func (config *TargetConfig) syscallOrInterrupt(st Statement, syscall bool) string {
	var i int

	if !syscall {
		// Remove st[1], if it's not a value
		i = 1
		if st[i].T != VALUE {
			//	log.Println("REMOVING ", st[i]);
			st = st[:i+copy(st[i:], st[i+1:])]
		}
	}

	// Remove st[-1] if it's a SEP
	i = len(st) - 1
	if st[i].T == SEP {
		st = st[:i+copy(st[i:], st[i+1:])]
	}

	log.Println("system call:")
	for _, token := range st {
		log.Println(token)
	}

	// Store each of the parameters to the appropriate registers
	var reg, n, comment, asmcode, precode, postcode string

	// How many tokens to skip before start reading arguments
	preskip := 2
	if syscall {
		preskip = 1
	}

	fromI := preskip //inclusive
	toI := len(st)   // exclusive
	stepI := 1
	if config.macOS {
		// arguments are pushed in the opposite order for BSD/OSX (32-bit)
		fromI = len(st) - 1 // inclusive
		toI = 1             // exclusive
		stepI = -1
	}
	firstI := fromI      // 2 for others, len(st)=1 for OSX/BSD
	lastI := toI - stepI // 2 for OSX/BSD, len(st)-1 for others
	for i := fromI; i != toI; i += stepI {
		if (i - preskip) >= len(config.interruptParameterRegisters) {
			log.Println("Error: Too many parameters for interrupt call:")
			for _, t := range st {
				log.Println(t.Value)
			}
			os.Exit(1)
			break
		}
		reg = config.interruptParameterRegisters[i-preskip]
		n = strconv.Itoa(i - preskip)
		if (config.macOS && (i == lastI)) || (!config.macOS && (i == firstI)) {
			comment = "function call: " + st[i].Value
		} else {
			if st[i].T == VALUE {
				comment = "parameter #" + n + " is " + st[i].Value
			} else {
				if strings.HasPrefix(st[i].Value, "_length_of_") {
					comment = "parameter #" + n + " is len(" + st[i].Value[11:] + ")"
				} else {
					if st[i].Value == "_" {
						// When _ is given, use the value already in the corresponding register
						comment = "parameter #" + n + " is supposedly already set"
					} else if has(dataNotValueTypes, st[i].Value) {
						comment = "parameter #" + n + " is " + "&" + st[i].Value
					} else {
						comment = "parameter #" + n + " is " + st[i].Value
						// Already recognized not to be a register
						switch config.PlatformBits {
						case 64:
							if st[i].Value == "rsp" {
								if is64bit(st[i].extra) {
									// Put the value of the register associated with this token at rbp
									precode += "\tsub rsp, 8\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
									precode += "\tmov QWORD [rsp], " + st[i].extra + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
									postcode += "\tadd rsp, 8\t\t\t; move the stack pointer back\n"
									break
								} else if is32bit(st[i].extra) {
									// Put the value of the register associated with this token at rbp
									precode += "\tsub rsp, 8\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
									precode += "\tmov QWORD [rsp], " + upgrade(st[i].extra) + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
									postcode += "\tadd rsp, 8\t\t\t; move the stack pointer back\n"
									break
								} else if is16bit(st[i].extra) {
									// Put the value of the register associated with this token at rbp
									precode += "\tsub rsp, 8\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
									precode += "\tmov QWORD [rsp], " + upgrade(upgrade(st[i].extra)) + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
									postcode += "\tadd rsp, 8\t\t\t; move the stack pointer back\n"
									break
								}
								log.Fatalln("Error: Unhandled register:", st[i].extra)
							}
						case 32:
							if st[i].Value == "esp" {
								if is32bit(st[i].extra) {
									precode += "\tsub esp, 4\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
									precode += "\tmov DWORD [esp], " + st[i].extra + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
									postcode += "\tadd esp, 4\t\t\t; move the stack pointer back\n"
									break
								} else if is16bit(st[i].extra) {
									precode += "\tsub esp, 4\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
									precode += "\tmov DWORD [esp], " + upgrade(st[i].extra) + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
									postcode += "\tadd esp, 4\t\t\t; move the stack pointer back\n"
									break
								}
								log.Fatalln("Error: Unhandled register:", st[i].extra)
							}
						case 16:
							// TODO: Add check for 8-bit values too: "mov BYTE [esp]"
							//log.Fatalln("Error: PARAMETERS are not implemented for 16-bit, yet")
							precode += "\tsub sp, 2\t\t\t; make some space for storing " + st[i].extra + " on the stack\n"
							precode += "\tmov WORD [sp], " + st[i].extra + "\t\t; move " + st[i].extra + " to a memory location on the stack\n"
							postcode += "\tadd sp, 2\t\t\t; move the stack pointer back\n"

						}
					}
				}
			}
		}
		codeline := ""
		// Skip parameters/registers that are already set
		if st[i].Value == "_" {
			codeline += "\t\t"

		} else {
			if st[i].Value == "0" {
				codeline += "\txor " + reg + ", " + reg
			} else {
				if config.macOS {
					if i == lastI {
						codeline += "\tmov " + reg + ", " + st[i].Value
					} else {
						codeline += "\tpush dword " + st[i].Value
					}
				} else {
					codeline += "\tmov " + reg + ", " + st[i].Value
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
		// TODO: comment which system call it is, ie "print"
		precode = "\t;--- system call ---\n" + precode
	} else {
		comment := "\t;--- call interrupt "
		if !strings.HasPrefix(st[1].Value, "0x") {
			// add 0x if missing, assume interrupts will always be called by hex
			comment += "0x"
		}
		comment += st[1].Value + " ---\n"
		precode = comment + precode
	}
	// Add the interrupt call
	if syscall || (st[1].T == VALUE) {
		if config.macOS {
			// just the way function calls are made on BSD/OSX
			asmcode += "\tsub esp, 4\t\t\t; BSD system call preparation\n"
		}
		if syscall {
			asmcode += "\tsyscall\t\t\t\t; perform the call\n"
		} else {
			// Add 0x if missing, assume interrupts will always be called by hex
			asmcode += "\tint "
			if !strings.HasPrefix(st[1].Value, "0x") {
				log.Println("Note: Adding 0x in front of interrupt", st[1].Value)
				asmcode += "0x"
			}
			asmcode += st[1].Value + "\t\t\t; perform the call\n"
		}
		if config.macOS {
			pushcount := len(st) - 2
			displacement := strconv.Itoa(pushcount * 4) // 4 bytes per push
			asmcode += "\tadd esp, " + displacement + "\t\t\t; BSD system call cleanup\n"
		}
		return precode + asmcode + postcode
	}
	log.Fatalln("Error: Need a (hexadecimal) interrupt number to call:\n", st[1].Value)
	return ""
}

func (st Statement) String(ps *ProgramState, config *TargetConfig) string {
	debug := true

	var parseState ParseState

	reduced := config.reduce(st, debug, ps)
	if len(reduced) != len(st) {
		return reduced.String(ps, config)
	}
	if len(st) == 0 {
		log.Fatalln("Error: Empty statement.")
		return ""
	} else if (st[0].T == BUILTIN) && (st[0].Value == "int") { // interrrupt call
		return config.syscallOrInterrupt(st, false)
	} else if (st[0].T == BUILTIN) && (st[0].Value == "syscall") {
		return config.syscallOrInterrupt(st, true)
	} else if (st[0].T == KEYWORD) && (st[0].Value == "var") && (len(st) >= 3) { // variable / bss declaration
		varname := ""
		if st[1].T == VALIDNAME {
			varname = st[1].Value
		} else {
			log.Fatalln("Error: "+st[1].Value, "is not a valid name for a variable")
		}
		bsscode := ""
		if (st[1].T == VALIDNAME) && ((st[2].T == VALUE) || (strings.HasPrefix(st[2].Value, "_length_of_"))) {
			if has(ps.definedNames, varname) {
				log.Fatalln("Error: Can not declare variable, name is already defined: " + varname)
			}
			ps.definedNames = append(ps.definedNames, varname)
			// Store the name of the declared variable in variables + the length
			if !strings.HasPrefix(st[2].Value, "_length_of_") {
				var err error
				ps.variables[varname], err = strconv.Atoi(st[2].Value)
				if err != nil {
					log.Fatalln("Error: " + st[2].Value + " is not a valid number of bytes to reserve")
				}
			}
			// Will be placed in the .bss section at the end
			bsscode += varname + ": resb " + st[2].Value + "\t\t\t\t; reserve " + st[2].Value + " bytes as " + varname + "\n"
			bsscode += "_capacity_of_" + varname + " equ " + st[2].Value + "\t\t; size of reserved memory\n"
			bsscode += "_length_of_" + varname + ": "
			switch config.PlatformBits {
			case 64:
				bsscode += "resd 1"
			case 32:
				bsscode += "resw 1"
			case 16:
				bsscode += "resb 1"
			}
			bsscode += "\t\t; current length of contents (points to after the data)\n"
			return bsscode
		}
		log.Printf("Error: Variable statements are on the form: \"var x 1024\" for reserving 1024 bytes, not: %s %s %s\n", st[0].Value, st[1].Value, st[2].Value)
		log.Println("Invalid parameters for variable string statement:")
		for _, t := range st {
			log.Println(t.Value)
		}
		os.Exit(1)
	} else if (st[0].T == KEYWORD) && (st[0].Value == "const") && (len(st) >= 4) { // constant data
		constname := ""
		if st[1].T == VALIDNAME {
			constname = st[1].Value
		} else {
			log.Fatalln("Error: "+st[1].Value, " (or a,b,c,d) is not a valid name for a constant")
		}
		asmcode := ""
		if (st[1].T == VALIDNAME) && (st[2].T == ASSIGNMENT) && ((st[3].T == STRING) || (st[3].T == VALUE) || (st[3].T == VALIDNAME)) {
			if has(ps.definedNames, constname) {
				log.Fatalln("Error: Can not declare constant, name is already defined: " + constname)
			}
			if (st[3].T == VALIDNAME) && !has(ps.definedNames, st[3].Value) {
				log.Fatalln("Error: Can't assign", st[3].Value, "to", st[1].Value, "because", st[3].Value, "is undefined.")
			}
			// Store the name of the declared constant in defined_names
			ps.definedNames = append(ps.definedNames, constname)
			// For the .DATA section (recognized by the keyword)
			if st[3].T == VALUE {
				switch config.PlatformBits {
				case 64:
					asmcode += constname + ":\tdq "
				case 32:
					asmcode += constname + ":\tdw "
				case 16:
					asmcode += constname + ":\tdb "
				}
			} else {
				asmcode += constname + ":\tdb "
				dataNotValueTypes = append(dataNotValueTypes, constname)
			}
			for i := 3; i < len(st); i++ {
				asmcode += st[i].Value
				// Add a comma between every element but the last one
				if (i + 1) != len(st) {
					asmcode += ", "
				}
			}
			if st[3].T == STRING {
				asmcode += "\t\t; constant string\n"
				//if config.platformBits == 16 {
				// Add an extra $, for safety, if on a 16-bit platform. Needed for print().
				// TODO: Remove, use a different int 21h call instead!
				//asmcode += "\tdb \"$\"\t\t\t; end of string, for when using ah=09/int 21h\n"
				//}
			} else {
				asmcode += "\t\t; constant value\n"
			}
			// Special naming for storing the length for later
			asmcode += "_length_of_" + constname + " equ $ - " + constname + "\t; size of constant value\n"
			return asmcode
		}
		log.Println("Error: Invalid parameters for constant string statement:")
		for _, t := range st {
			log.Println(t.Value)
		}
		os.Exit(1)
	} else if (len(st) > 2) && (st[0].T == VALIDNAME) && (st[1].T == ASSIGNMENT) {
		// Copying data from constants to variables (reserved memory in the .bss section)
		asmcode := ""
		from := st[2].Value
		to := st[0].Value
		lengthexpr := "_length_of_" + from
		toPosition := "[_length_of_" + to + "]"
		// TODO: Make this a lot smarter and handle copying ranges of data, adr or value
		// TODO: Actually, redesign the whole language
		switch config.PlatformBits {
		case 64:
			asmcode += "\tmov rdi, " + to + "\t\t\t; copy bytes from " + from + " to " + to + "\n"
			asmcode += "\tmov rsi, " + from + "\n"
			asmcode += "\tmov rcx, " + lengthexpr + "\n"
			//asmcode += "\tmov QWORD " + toPosition + ", " + to + "\n"
			asmcode += "\tmov " + toPosition + ", rcx" + "\n"
			asmcode += "\tcld\n"
			asmcode += "\trep movsb\t\t\t\t; copy bytes\n" // optimized ok on 64-bit CPUs
		case 32:
			asmcode += "\tmov edi, " + to + "\t\t\t; copy bytes from " + from + " to " + to + "\n"
			asmcode += "\tmov esi, " + from + "\n"
			asmcode += "\tmov ecx, " + lengthexpr + "\n"
			asmcode += "\tmov " + toPosition + ", ecx\n"
			asmcode += "\tcld\n"
			asmcode += "\trep movsb\t\t\t\t; copy bytes\n" // optimized ok on 32-bit CPUs
		case 16:
			// TODO: Test this
			asmcode += "\tmov di, " + to + "\t\t\t; copy bytes from " + from + " to " + to + "\n"
			asmcode += "\tmov si, " + from + "\n"
			asmcode += "\tmov cx, " + lengthexpr + "\n"
			asmcode += "\tmov " + toPosition + ", cx\n"
			asmcode += "\trep movsb\t\t\t\t; copy bytes\n"
		}
		return asmcode
	} else if (len(st) > 2) && ((st[1].T == ADDITION) && (st[0].T == VALIDNAME) && (st[2].T == VALIDNAME)) {
		// Copying data from constants to variables (reserved memory in the .bss section)
		asmcode := ""
		from := st[2].Value
		to := st[0].Value
		lengthAddr := "[_length_of_" + to + "]"
		// TODO: Make this a lot smarter and handle copying ranges of data, adr or value
		// TODO: Actually, redesign the whole language
		switch config.PlatformBits {
		case 64:
			asmcode += "\tmov rdi, " + to + "\t\t; add bytes from \"" + from + "\" to " + to + "\n"
			asmcode += "\tadd rdi, " + lengthAddr + "\n"
			asmcode += "\tmov rsi, " + from + "\n"
			asmcode += "\tmov rcx, _length_of_" + from + "\n"
			asmcode += "\tadd " + lengthAddr + ", rcx" + "\n"
			asmcode += "\tcld\n"
			asmcode += "\trep movsb\t\t\t\t; copy bytes\n"
		case 32:
			asmcode += "\tmov edi, " + to + "\t\t; add bytes from \"" + from + "\" to " + to + "\n"
			asmcode += "\tadd edi, " + lengthAddr + "\n"
			asmcode += "\tmov esi, " + from + "\n"
			asmcode += "\tmov ecx, _length_of_" + from + "\n"
			asmcode += "\tadd " + lengthAddr + ", ecx" + "\n"
			asmcode += "\tcld\n"
			asmcode += "\trep movsb\t\t\t\t; copy bytes\n"
		case 16:
			// TODO: Test this
			asmcode += "\tmov di, " + to + "\t\t; add bytes from \"" + from + "\" to " + to + "\n"
			asmcode += "\tadd di, " + lengthAddr + "\n"
			asmcode += "\tmov si, " + from + "\n"
			asmcode += "\tmov cx, _length_of_" + from + "\n"
			asmcode += "\tadd " + lengthAddr + ", cx" + "\n"
			asmcode += "\trep movsb\t\t\t\t; copy bytes\n"
		}
		return asmcode
	} else if (st[0].T == BUILTIN) && (st[0].Value == "halt") {
		asmcode := "\t; --- full stop ---\n"
		asmcode += "\tcli\t\t; clear interrupts\n"
		asmcode += ".hang:\n"
		asmcode += "\thlt\n"
		asmcode += "\tjmp .hang\t; loop forever\n\n"
		return asmcode
	} else if (config.PlatformBits == 16) && (st[0].T == BUILTIN) && (st[0].Value == "print") && (st[1].T == VALIDNAME) {
		asmcode := "\t; --- output string of given length ---\n"
		asmcode += "\tmov dx, " + st[1].Value + "\n"
		if _, ok := ps.variables[st[1].Value]; ok {
			// A variable in .bss
			asmcode += "\tmov cx, [_length_of_" + st[1].Value + "]\n"
		} else {
			asmcode += "\tmov cx, _length_of_" + st[1].Value + "\n"
		}
		asmcode += "\tmov bx, 1\n"
		asmcode += "\tmov ah, 0x40\t\t; prepare to call \"Write File or Device\"\n"
		asmcode += "\tint 0x21\n\n"
		return asmcode
	} else if ((st[0].T == KEYWORD) && (st[0].Value == "ret")) || ((st[0].T == BUILTIN) && (st[0].Value == "exit")) {
		asmcode := ""
		if st[0].Value == "ret" {
			if (ps.inFunction == "main") || (ps.inFunction == config.LinkerStartFunction) {
				//log.Println("Not taking down stack frame in the main/_start/start function.")
			} else {
				switch config.PlatformBits {
				case 64:
					asmcode += "\t;--- takedown stack frame ---\n"
					asmcode += "\tmov rsp, rbp\t\t\t; use base pointer as new stack pointer\n"
					asmcode += "\tpop rbp\t\t\t\t; get the old base pointer\n\n"
				case 32:
					asmcode += "\t;--- takedown stack frame ---\n"
					asmcode += "\tmov esp, ebp\t\t\t; use base pointer as new stack pointer\n"
					asmcode += "\tpop ebp\t\t\t\t; get the old base pointer\n\n"
				}
			}
		}
		if ps.inFunction != "" {
			if !config.BootableKernel && !ps.endless && (ps.inFunction == "main") {
				asmcode += "\n\t;--- return from \"" + ps.inFunction + "\" ---\n"
			}
		} else if st[0].Value == "exit" {
			asmcode += "\t;--- exit program ---\n"
		} else {
			asmcode += "\t;--- return ---\n"
		}
		if (st[0].Value == "exit") || (ps.inFunction == "main") || (ps.inFunction == config.LinkerStartFunction) {
			// Not returning from main/_start/start function, but exiting properly
			exitCode := "0"
			if (len(st) == 2) && ((st[1].T == VALUE) || (st[1].T == REGISTER)) {
				exitCode = st[1].Value
			}
			if !config.BootableKernel {
				switch config.PlatformBits {
				case 64:
					asmcode += "\tmov rax, 60\t\t\t; function call: 60\n\t"
					if exitCode == "0" {
						asmcode += "xor rdi, rdi"
					} else {
						asmcode += "mov rdi, " + exitCode
					}
					asmcode += "\t\t\t; return code " + exitCode + "\n"
					asmcode += "\tsyscall\t\t\t\t; exit program\n"
				case 32:
					if config.macOS {
						asmcode += "\tpush dword " + exitCode + "\t\t\t; exit code " + exitCode + "\n"
						asmcode += "\tsub esp, 4\t\t\t; the BSD way, push then subtract before calling\n"
					}
					asmcode += "\tmov eax, 1\t\t\t; function call: 1\n"
					if !config.macOS {
						asmcode += "\t"
						if exitCode == "0" {
							asmcode += "xor ebx, ebx"
						} else {
							asmcode += "mov ebx, " + exitCode
						}
						asmcode += "\t\t\t; exit code " + exitCode + "\n"
					}
					asmcode += "\tint 0x80\t\t\t; exit program\n"
				case 16:
					// Unless "exit" or "noret" is specified explicitly, use "ret"
					if st[0].Value == "exit" {
						// Since we are not building a kernel, calling DOS interrupt 21h makes sense
						asmcode += "\tmov ah, 0x4c\t\t\t; function 4C\n"
						if exitCode == "0" {
							asmcode += "\txor al, al\t\t\t; exit code " + exitCode + "\n"
						} else {
							asmcode += "\tmov al, " + exitCode + "\t\t\t; exit code " + exitCode + "\n"
						}
						asmcode += "\tint 0x21\t\t\t; exit program\n"
					} else if st[0].Value == "noret" {
						asmcode += "\t; there is no return\n"
					} else {
						if !ps.endless {
							asmcode += "\tret\t\t\t; exit program\n"
						} else {
							asmcode += "\t; endless loop, there is no return\n"
						}
					}
				}
			} else {
				// For bootable kernels, main does not return. Hang instead.
				log.Println("Warning: Bootable kernels has nowhere to return after the main function. You might want to use the \"halt\" builtin at the end of the main function.")
				//asmcode += Statement{Token{BUILTIN, "halt", st[0].line, ""}}.String()
			}
		} else {
			log.Println("function ", ps.inFunction)
			// Do not return eax=0/rax=0 if no return value is explicitly provided, by design
			// This allows the return value from the previous call to be returned instead
			asmcode += "\tret\t\t\t\t; Return\n"
		}
		if ps.inFunction != "" {
			// Exiting from the function definition
			ps.inFunction = ""
			// If the function was ended with "exit", don't freak out if an "end" is encountered
			if st[0].Value == "exit" {
				ps.surpriseEndingWithExit = true
			}
		}
		if parseState.inlineC {
			// Exiting from inline C
			parseState.inlineC = false
			return "; End of inline C block"
		}
		return asmcode
	} else if (st[0].T == KEYWORD && st[0].Value == "mem") && (st[1].T == VALUE || st[1].T == VALIDNAME || st[1].T == REGISTER) && (st[2].T == ASSIGNMENT) && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// memory assignment
		return "\tmov [" + st[1].Value + "], " + st[3].Value + "\t\t; " + "memory assignment" + "\n"
	} else if (st[0].T == KEYWORD && st[0].Value == "membyte") && (st[1].T == VALUE || st[1].T == VALIDNAME || st[1].T == REGISTER) && (st[2].T == ASSIGNMENT) && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// memory assignment (byte)
		val := st[3].Value
		if st[3].T == REGISTER {
			val = downgradeToByte(val)
		}
		return "\tmov BYTE [" + st[1].Value + "], " + val + "\t\t; " + "memory assignment" + "\n"
	} else if (st[0].T == KEYWORD && st[0].Value == "memword") && (st[1].T == VALUE || st[1].T == VALIDNAME || st[1].T == REGISTER) && (st[2].T == ASSIGNMENT) && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// memory assignment (word)
		val := st[3].Value
		if st[3].T == REGISTER {
			val = regToWord(val)
		}
		return "\tmov WORD [" + st[1].Value + "], " + val + "\t\t; " + "memory assignment" + "\n"
	} else if (st[0].T == KEYWORD && st[0].Value == "memdouble") && (st[1].T == VALUE || st[1].T == VALIDNAME || st[1].T == REGISTER) && (st[2].T == ASSIGNMENT) && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// memory assignment (double)
		val := st[3].Value
		if st[3].T == REGISTER {
			val = regToDouble(val)
		}
		return "\tmov DOUBLE [" + st[1].Value + "], " + val + "\t\t; " + "memory assignment" + "\n"
	} else if (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == KEYWORD && st[2].Value == "mem") && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// assignment from memory to register
		return "\tmov " + st[0].Value + ", [" + st[3].Value + "]\t\t; memory assignment\n"
	} else if (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == KEYWORD && st[2].Value == "readbyte") && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// assignment from memory to register (byte)
		val := st[0].Value
		if st[0].T == REGISTER {
			val = downgradeToByte(val)
		}
		return "\tmov BYTE " + val + ", [" + st[3].Value + "]\t\t; memory assignment (byte)\n"
	} else if (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == KEYWORD && st[2].Value == "readword") && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// assignment from memory to register (byte)
		val := st[0].Value
		if st[0].T == REGISTER {
			val = regToWord(val)
		}
		return "\tmov WORD " + val + ", [" + st[3].Value + "]\t\t; memory assignment (word)\n"
	} else if (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == KEYWORD && st[2].Value == "readdouble") && (st[3].T == VALUE || st[3].T == VALIDNAME || st[3].T == REGISTER) {
		// assignment from memory to register (byte)
		val := st[0].Value
		if st[0].T == REGISTER {
			val = regToDouble(val)
		}
		return "\tmov DOUBLE " + val + ", [" + st[3].Value + "]\t\t; memory assignment (double)\n"
	} else if len(st) == 3 && ((st[0].T == REGISTER) || (st[0].T == DISREGARD) || (st[0].Value == "stack") || (st[2].Value == "stack")) {
		// Statements like "eax = 3" are handled here
		// TODO: Handle all sorts of equivivalents to assembly statements
		if st[1].T == COMPARISON {
			if ps.inIfBlock != "" {
				log.Fatalln("Error: Already in an if-block (nested block are to be implemented)")
			}
			ps.inIfBlock = ps.newIfLabel()

			asmcode := "\t;--- " + ps.inIfBlock + " ---\n"

			// Start an if block that is run if the comparison is true
			// Break if something comparison something
			asmcode += "\tcmp " + st[0].Value + ", " + st[2].Value + "\t\t\t; compare\n"

			// Conditional jump if NOT true
			asmcode += "\t"
			switch st[1].Value {
			case "==":
				asmcode += "jne"
			case "!=":
				asmcode += "je"
			case ">":
				asmcode += "jle"
			case "<":
				asmcode += "jge"
			case "<=":
				asmcode += "jg"
			case ">=":
				asmcode += "jl"
			}

			// Which label to jump to (out of the if block)
			// TODO: Nested if blocks
			asmcode += " " + ps.inIfBlock + "_end\t\t\t; break\n"
			return asmcode
		} else if (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == VALUE || st[2].T == VALIDNAME) {
			if st[2].Value == "0" {
				return "\txor " + st[0].Value + ", " + st[0].Value + "\t\t; " + st[0].Value + " " + st[1].Value + " " + st[2].Value
			}
			a := st[0].Value
			b := st[2].Value
			if is32bit(a) && is64bit(b) {
				log.Println("Warning: Using", b, "as a 32-bit register when assigning.")
				return "\tmov " + a + ", " + downgrade(b) + "\t\t; " + a + " " + st[1].Value + " " + b
			} else if is64bit(a) && is32bit(b) {
				log.Println("Warning: Using", a, "as a 32-bit register when assigning.")
				asmcode := "\txor rax, rax\t\t; clear rax\n"
				asmcode += "\tmov " + downgrade(a) + ", " + b + "\t\t; " + a + " " + st[1].Value + " " + b
				return asmcode
			} else {
				return "\tmov " + st[0].Value + ", " + st[2].Value + "\t\t; " + st[0].Value + " " + st[1].Value + " " + st[2].Value
			}
		} else if (st[0].T == VALIDNAME) && (st[1].T == ASSIGNMENT) {
			if has(ps.definedNames, st[0].Value) {
				log.Fatalln("Error:", st[0].Value, "has already been defined")
			} else {
				log.Fatalln("Error:", st[0].Value, "is not recognized as a register (and there is no const qualifier). Can't assign.")
			}
		} else if st[0].T == DISREGARD {
			// TODO: If st[2] is a function, one wishes to call it, then disregard afterwards
			return "\t\t\t\t; Disregarding: " + st[2].Value + "\n"
		} else if (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == REGISTER) {
			return "\tmov " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " " + st[1].Value + " " + st[2].Value
		} else if (st[0].T == RESERVED) && (st[1].T == VALUE) {
			return config.reservedAndValue(st[:2])
		} else if (len(st) == 3) && ((st[0].T == REGISTER) || (st[0].Value == "stack") || (st[0].T == VALUE)) && (st[1].T == ARROW) && ((st[2].T == REGISTER) || (st[2].Value == "stack")) {
			// push and pop
			if (st[0].Value == "stack") && (st[2].Value == "stack") {
				log.Fatalln("Error: can't pop and push to stack at the same time")
			} else if st[2].Value == "stack" {
				// something -> stack (push)
				return "\tpush " + st[0].Value + "\t\t\t; " + st[0].Value + " -> stack\n"
			} else if st[0].Value == "stack" {
				// stack -> something (pop)
				return "\tpop " + st[2].Value + "\t\t\t\t; stack -> " + st[2].Value + "\n"
			} else if (st[0].T == REGISTER) && (st[2].T == REGISTER) {
				// reg -> reg (push and then pop)
				return "\tpush " + st[0].Value + "\t\t\t; " + st[0].Value + " -> " + st[2].Value + "\n\tpop " + st[2].Value + "\t\t\t\t;\n"
			} else {
				log.Println("Error: Unrecognized stack expression: ")
				for _, token := range []Token(st) {
					log.Print(token)
				}
				os.Exit(1)
			}
		} else if (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == RESERVED || st[2].T == VALUE) && (st[3].T == VALUE) {
			if st[2].Value == "funparam" {
				paramoffset, err := strconv.Atoi(st[3].Value)
				if err != nil {
					log.Fatalln("Error: Invalid list offset for", st[2].Value+":", st[3].Value)
				}
				paramExpression := config.paramnum2reg(paramoffset)
				if len(paramExpression) == 3 {
					paramExpression += "\t"
				}
				return "\tmov " + st[0].Value + ", " + paramExpression + "\t\t; fetch function param #" + st[3].Value + "\n"
			}
			// TODO: Implement support for other lists
			log.Fatalln("Error: Can only handle \"funparam\" lists when assigning to a register, so far.")
		}
		if (st[1].T == ADDITION) && (st[2].T == REGISTER) {
			return "\tadd " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " += " + st[2].Value
		} else if (st[1].T == SUBTRACTION) && (st[2].T == REGISTER) {
			return "\tsub " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " -= " + st[2].Value
		} else if (st[1].T == MULTIPLICATION) && (st[2].T == REGISTER) {
			if registerA(st[0].Value) {
				return "\tmul " + st[2].Value + "\t\t\t; " + st[0].Value + " *= " + st[2].Value
			}
			if st[0].Value == st[2].Value {
				return "\timul " + st[0].Value + "\t\t\t; " + st[0].Value + " *= " + st[0].Value
			}
			return "\timul " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " *= " + st[2].Value
		} else if (st[1].T == DIVISION) && (st[2].T == REGISTER) {
			if registerA(st[0].Value) {
				return "\tdiv " + st[2].Value + "\t\t\t; " + st[0].Value + " /= " + st[2].Value
			}
			return "\tidiv " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " /= " + st[2].Value
		}
		if (st[1].T == ADDITION) && ((st[2].T == VALUE) || (st[2].T == MEMEXP)) {
			if st[2].Value == "1" {
				return "\tinc " + st[0].Value + "\t\t\t; " + st[0].Value + "++"
			}
			return "\tadd " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " += " + st[2].Value
		} else if (st[1].T == SUBTRACTION) && ((st[2].T == VALUE) || (st[2].T == MEMEXP)) {
			if st[2].Value == "1" {
				return "\tdec " + st[0].Value + "\t\t\t; " + st[0].Value + "--"
			}
			return "\tsub " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " -= " + st[2].Value
		} else if (st[1].T == AND) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\tand " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " &= " + st[2].Value
		} else if (st[1].T == OR) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\tor " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " |= " + st[2].Value
			// TODO: All == MEMEXP should be followed by || st[2].t == REGEXP. In fact,
			//       a better system is needed. Some sort of pattern matching.
		} else if (st[1].T == XOR) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\txor " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " ^= " + st[2].Value
		} else if (st[1].T == ROL) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\trol " + st[0].Value + ", " + st[2].Value + "\t\t\t; rotate " + st[0].Value + " left" + st[2].Value
		} else if (st[1].T == ROR) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\tror " + st[0].Value + ", " + st[2].Value + "\t\t\t; rotate " + st[0].Value + " right " + st[2].Value
		} else if (st[1].T == SHL) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\tshl " + st[0].Value + ", " + st[2].Value + "\t\t\t; shift " + st[0].Value + " left" + st[2].Value
		} else if (st[1].T == SHR) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\tshr " + st[0].Value + ", " + st[2].Value + "\t\t\t; shift " + st[0].Value + " right " + st[2].Value
		} else if (st[1].T == XCHG) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\txchg " + st[0].Value + ", " + st[2].Value + "\t\t\t; exchange " + st[0].Value + " and " + st[2].Value
		} else if (st[1].T == OUT) && ((st[2].T == VALUE) || (st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\tout " + st[0].Value + ", " + st[2].Value + "\t\t\t; output " + st[0].Value + " to IO port " + st[2].Value
		} else if (st[1].T == IN) && ((st[2].T == MEMEXP) || (st[2].T == REGISTER)) {
			return "\tin " + st[2].Value + ", " + st[0].Value + "\t\t\t; input " + st[2].Value + " from IO port " + st[0].Value
		} else if (st[1].T == MULTIPLICATION) && ((st[2].T == VALUE) || (st[2].T == MEMEXP)) {
			// TODO: Don't use a list, write a function that covers the lot
			shifts := []string{"2", "4", "8", "16", "32", "64", "128"}
			if has(shifts, st[2].Value) {
				pos := 0
				for i, v := range shifts {
					if v == st[2].Value {
						// Found the appropriate shift value
						pos = i + 1
						break
					}
				}
				// TODO: Check that it works with signed numbers and/or introduce signed/unsigned operations
				return "\tshl " + st[0].Value + ", " + strconv.Itoa(pos) + "\t\t\t; " + st[0].Value + " *= " + st[2].Value
			}
			if registerA(st[0].Value) {
				return "\tmul " + st[2].Value + "\t\t\t; " + st[0].Value + " *= " + st[2].Value
			}
			if st[0].Value == st[2].Value {
				return "\timul " + st[0].Value + "\t\t\t; " + st[0].Value + " *= " + st[0].Value
			}
			return "\timul " + st[0].Value + ", " + st[2].Value + "\t\t\t; " + st[0].Value + " *= " + st[2].Value
		} else if (st[1].T == DIVISION) && ((st[2].T == VALUE) || (st[2].T == MEMEXP)) {
			// TODO: Don't use a list, write a function that covers the lot
			shifts := []string{"2", "4", "8", "16", "32", "64", "128"}
			if has(shifts, st[2].Value) {
				pos := 0
				for i, v := range shifts {
					if v == st[2].Value {
						// Found the appropriate shift value
						pos = i + 1
						break
					}
				}
				// TODO: Check that it works with signed numbers and/or introduce signed/unsigned operations
				return "\tshr " + st[0].Value + ", " + strconv.Itoa(pos) + "\t\t; " + st[0].Value + " /= " + st[2].Value
			}
			asmcode := "\n\t;--- signed division: " + st[0].Value + " /= " + st[2].Value + " ---\n"
			// TODO Add support for division with 16-bit registers as well!

			if config.PlatformBits == 32 {
				if st[0].Value == "eax" {
					// Dividing a 64-bit number in edx:eax by the number in ecx. Clearing out edx and only using 32-bit numbers for now.
					// If the register to be divided is rax, do a quicker division than if it's another register

					// save ecx
					asmcode += "\tpush ecx\t\t; save ecx\n"
					//// save edx
					//asmcode += "\tpush edx\t\t; save edx\n"
					// clear edx
					asmcode += "\txor edx, edx\t\t; edx = 0 (32-bit 0:eax instead of 64-bit edx:eax)\n"
					// ecx = st[2].value
					asmcode += "\tmov ecx, " + st[2].Value + "\t\t; divisor, ecx = " + st[2].Value + "\n"
					// div ecx
					asmcode += "\tdiv ecx\t\t\t; eax = edx:eax / ecx\n"
					asmcode += "\t\t\t; remainder is in edx\n"
					//// restore edx
					//asmcode += "\tpop edx\t\t; restore edx\n"
					// restore ecx
					asmcode += "\tpop ecx\t\t; restore ecx\n"
				} else if st[0].Value == "ax" {
					// Dividing a 32-bit number in dx:ax by the number in bx. Clearing out dx and only using 16-bit numbers for now.
					// If the register to be divided is ax, do a quicker division than if it's another register

					// save bx
					asmcode += "\tpush cx\t\t; save cx\n"
					//// save dx
					//asmcode += "\tpush dx\t\t; save dx\n"
					// clear dx
					asmcode += "\txor dx, dx\t; dx = 0 (16-bit 0:ax instead of 32-bit dx:ax)\n"
					// bx = st[2].value
					asmcode += "\tmov cx, " + st[2].Value + "\t; divisor, cx = " + st[2].Value + "\n"
					asmcode += "\t\t\t; remainder is in dx\n"
					// div bx
					asmcode += "\tdiv cx\t\t; ax = dx:ax / cx\n"
					//// restore dx
					//asmcode += "\tpop dx\t\t; restore dx\n"
					// restore cx
					asmcode += "\tpop cx\t\t; restore cx\n"
				} else {
					// TODO: if the given register is a different one than eax, ecx and edx,
					//       just divide directly with that register, like for eax above
					// save eax, we know this is not where we assign the result
					asmcode += "\tpush eax\t\t; save eax\n"
					if st[0].Value != "ecx" {
						// save ecx
						asmcode += "\tpush ecx\t\t; save ecx\n"
					}
					if st[0].Value != "edx" {
						// save edx
						asmcode += "\tpush edx\t\t; save edx\n"
					}
					// copy number to be divided to eax
					if is64bit(st[0].Value) {
						if downgrade(st[0].Value) != "eax" {
							asmcode += "\tmov eax, " + downgrade(st[0].Value) + "\t\t; dividend, number to be divided\n"
						}
					} else if is16bit(st[0].Value) {
						if upgrade(st[0].Value) != "eax" {
							asmcode += "\tmov eax, " + upgrade(st[0].Value) + "\t\t; dividend, number to be divided\n"
						}
					} else {
						if st[0].Value != "eax" {
							asmcode += "\tmov eax, " + st[0].Value + "\t\t; dividend, number to be divided\n"
						}
					}
					// clear edx
					asmcode += "\txor edx, edx\t\t; edx = 0 (32-bit 0:eax instead of 64-bit edx:eax)\n"
					// ecx = st[2].value
					asmcode += "\tmov ecx, " + st[2].Value + "\t\t; divisor, ecx = " + st[2].Value + "\n"
					// eax = edx:eax / ecx
					asmcode += "\tdiv ecx\t\t\t; eax = edx:eax / ecx\n"
					if st[0].Value != "edx" {
						// restore edx
						asmcode += "\tpop edx\t\t; restore edx\n"
					}
					if st[0].Value != "ecx" {
						// restore ecx
						asmcode += "\tpop ecx\t\t; restore ecx\n"
					}
					// st[0].value = eax
					asmcode += "\tmov " + st[0].Value + ", eax\t\t; " + st[0].Value + " = eax\n"
					// restore eax
					asmcode += "\tpop eax\t\t; restore eax\n"
				}
				asmcode += "\n"
				return asmcode
			}
			// Dividing a 128-bit number in rdx:rax by the number in rcx. Clearing out rdx and only using 64-bit numbers for now.
			// If the register to be divided is rax, do a quicker division than if it's another register
			if st[0].Value == "rax" {
				// save rdx
				//asmcode += "\tmov r9, rdx\t\t; save rdx\n"
				// clear rdx
				asmcode += "\txor rdx, rdx\t\t; rdx = 0 (64-bit 0:rax instead of 128-bit rdx:rax)\n"
				// mov r8, st[2].value
				asmcode += "\tmov r8, " + st[2].Value + "\t\t; divisor, r8 = " + st[2].Value + "\n"
				// div rax
				asmcode += "\tdiv r8\t\t\t; rax = rdx:rax / r8\n"
				// restore rdx
				//asmcode += "\tmov rdx, r9\t\t; restore rdx\n"
			} else {
				log.Println("Note: r8, r9 and r10 will be changed when dividing: " + st[0].Value + " /= " + st[2].Value)
				// TODO: if the given register is a different one than rax, rcx and rdx,
				//       just divide directly with that register, like for rax above
				// save rax, we know this is not where we assign the result
				if !registerA(st[0].Value) {
					asmcode += "\tmov r9, rax\t\t; save rax\n"
				}
				//if st[0].value != "rdx" {
				//	// save rdx
				//	asmcode += "\tmov r10, rdx\t\t; save rdx\n"
				//}
				// copy number to be divided to rax
				if is32bit(st[0].Value) {
					if st[0].Value != "eax" {
						asmcode += "\txor rax, rax\t\t; clear rax\n"
						asmcode += "\tmov eax, " + st[0].Value + "\t\t; dividend, number to be divided\n"
					}
				} else if is16bit(st[0].Value) {
					if st[0].Value != "ax" {
						asmcode += "\txor rax, rax\t\t; clear rax\n"
						asmcode += "\tmov ax, " + st[0].Value + "\t\t; dividend, number to be divided\n"
					}
				} else {
					if st[0].Value != "rax" {
						asmcode += "\tmov rax, " + st[0].Value + "\t\t; dividend, number to be divided\n"
					}
				}
				// xor rdx, rdx
				asmcode += "\txor rdx, rdx\t\t; rdx = 0 (64-bit 0:rax instead of 128-bit rdx:rax)\n"
				// mov rcx, st[2].value
				asmcode += "\tmov r8, " + st[2].Value + "\t\t; divisor, r8 = " + st[2].Value + "\n"
				// div rax
				asmcode += "\tdiv r8\t\t\t; rax = rdx:rax / r8\n"
				//if st[0].value != "rdx" {
				//	// restore rdx
				//	asmcode += "\tmov rdx, r10\t\t; restore rdx\n"
				//}
				// mov st[0].value, rax
				if !registerA(st[0].Value) {
					asmcode += "\tmov " + st[0].Value + ", rax\t\t; " + st[0].Value + " = rax\n"
				}
				// restore rax
				if !registerA(st[0].Value) {
					asmcode += "\tmov rax, r9\t\t; restore rax\n"
				}
			}
			return asmcode
		}
		log.Println("Unfamiliar 3-token expression!")
	} else if (len(st) == 4) && (st[0].T == RESERVED) && (st[1].T == VALUE) && (st[2].T == ASSIGNMENT) && ((st[3].T == VALIDNAME) || (st[3].T == VALUE) || (st[3].T == REGISTER)) {
		retval := "\tmov " + config.reservedAndValue(st[:2]) + ", " + st[3].Value + "\t\t\t; "
		if (config.PlatformBits == 32) && (st[3].T != REGISTER) {
			retval = strings.Replace(retval, "mov", "mov DWORD", 1)
		}
		pointercomment := ""
		if st[3].T == VALIDNAME {
			pointercomment = "&"
		}
		retval += fmt.Sprintf("%s[%s] = %s%s\n", st[0].Value, st[1].Value, pointercomment, st[3].Value)
		return retval
	} else if (len(st) == 4) && (st[0].T == REGISTER) && (st[1].T == ASSIGNMENT) && (st[2].T == RESERVED) && (st[3].T == VALUE) {
		retval := "\tmov " + st[0].Value + ", " + config.reservedAndValue(st[2:]) + "\t\t\t; "
		retval += fmt.Sprintf("%s = %s[%s]\n", st[0].Value, st[2].Value, st[3].Value)
		return retval
	} else if (len(st) == 5) && (st[0].T == RESERVED) && (st[1].T == VALUE) && (st[2].T == ASSIGNMENT) && (st[3].T == RESERVED) && (st[4].T == VALUE) {
		retval := ""
		if config.PlatformBits != 32 {
			retval = "\tmov " + config.reservedAndValue(st[:2]) + ", " + config.reservedAndValue(st[3:]) + "\t\t\t; "
		} else {
			retval = "\tmov eax, " + config.reservedAndValue(st[3:]) + "\t\t\t; Uses eax as a temporary variable\n"
			retval += "\tmov " + config.reservedAndValue(st[:2]) + ", ebx\t\t\t; "
		}
		retval += fmt.Sprintf("%s[%s] = %s[%s]\n", st[0].Value, st[1].Value, st[3].Value, st[4].Value)
		return retval
	} else if (len(st) >= 2) && (st[0].T == KEYWORD) && (st[0].Value == "asm") && (st[1].T == VALUE) {
		targetBits, err := strconv.Atoi(st[1].Value)
		if err != nil {
			log.Fatalln("Error: " + st[1].Value + " is not a valid platform bit size (like 32 or 64)")
		}
		if config.PlatformBits == targetBits {
			// Add the rest of the line as a regular assembly expression
			if len(st) == 7 {
				comma1 := " "
				comma2 := ", "
				if st[4].T == QUAL {
					comma1 = ", "
					comma2 = " "
				}
				// with address calculations
				if strings.Contains(st[5].Value, "+") || strings.Contains(st[5].Value, "-") {
					return "\t" + st[2].Value + " " + st[3].Value + " " + st[4].Value + " " + st[5].Value + " " + st[6].Value + "\t\t\t; asm with address calculation\n"
				} else if strings.HasPrefix(st[2].Value, "i") {
					comma1 = ", "
					return "\t" + st[2].Value + " " + st[3].Value + comma1 + st[4].Value + comma2 + st[5].Value + " " + st[6].Value + "\t\t\t; asm with integer maths\n"
				} else {
					return "\t" + st[2].Value + " " + st[3].Value + comma1 + st[4].Value + comma2 + st[5].Value + " " + st[6].Value + "\t\t\t; asm with floating point instructions\n"
				}
			} else if len(st) == 6 {
				comma1 := " "
				comma2 := ", "
				if st[4].T == QUAL {
					comma1 = ", "
					comma2 = " "
				}
				// with address calculations
				if strings.Contains(st[5].Value, "+") || strings.Contains(st[5].Value, "-") {
					return "\t" + st[2].Value + " " + st[3].Value + comma1 + st[4].Value + comma2 + st[5].Value + "\t\t\t; asm with address calculation\n"
				} else if strings.HasPrefix(st[2].Value, "i") {
					comma1 = ", "
					return "\t" + st[2].Value + " " + st[3].Value + comma1 + st[4].Value + comma2 + st[5].Value + "\t\t\t; asm with integer maths\n"
				} else {
					return "\t" + st[2].Value + " " + st[3].Value + comma1 + st[4].Value + comma2 + st[5].Value + "\t\t\t; asm with floating point instructions\n"
				}
			} else if len(st) == 5 {
				comma2 := ", "
				if st[3].T == QUAL {
					comma2 = " "
				}
				// with address calculations
				if strings.Contains(st[4].Value, "+") || strings.Contains(st[4].Value, "-") {
					return "\t" + st[2].Value + " " + st[3].Value + comma2 + st[4].Value + "\t\t\t; asm with address calculation\n"
				} else if st[3].Value == "st" {
					return "\t" + st[2].Value + " " + st[3].Value + " (" + st[4].Value + ")\t\t\t; asm\n"
				} else {
					return "\t" + st[2].Value + " " + st[3].Value + comma2 + st[4].Value + "\t\t\t; asm\n"
				}
			} else if len(st) == 4 {
				return "\t" + st[2].Value + " " + st[3].Value + "\t\t\t; asm\n"
			} else if len(st) == 3 {
				// a label or keyword like "stosb"
				if strings.Contains(st[2].Value, ":") {
					return "\t" + st[2].Value + "\t\t\t; asm label\n"
				}
				return "\t" + st[2].Value + "\t\t\t; asm\n"
			} else {
				log.Println("Error: Unrecognized length of assembly expression:", len(st)-2)
				for i, token := range []Token(st) {
					if i < 2 {
						continue
					}
					log.Print(token)
				}
				os.Exit(1)
			}
		}
		// Not the target bits, skip
		return ""
	} else if (len(st) >= 2) && (st[0].T == KEYWORD) && (st[1].T == VALIDNAME) && (st[0].Value == "fun") {
		if ps.inFunction != "" {
			log.Fatalf("Error: Missing \"ret\" or \"end\"? Already in a function named %s when declaring function %s.\n", ps.inFunction, st[1].Value)
		}
		asmcode := ";--- function " + st[1].Value + " ---\n"
		ps.inFunction = st[1].Value
		// Store the name of the declared function in defined_names
		if has(ps.definedNames, ps.inFunction) {
			log.Fatalln("Error: Can not declare function, name is already defined:", ps.inFunction)
		}
		ps.definedNames = append(ps.definedNames, ps.inFunction)
		if config.PlatformBits != 16 {
			asmcode += "global " + ps.inFunction + "\t\t\t; make label available to the linker\n"
		}
		asmcode += ps.inFunction + ":\t\t\t\t; name of the function\n\n"
		if (ps.inFunction == "main") || (ps.inFunction == config.LinkerStartFunction) {
			//log.Println("Not setting up stack frame in the main/_start/start function.")
			return asmcode
		}
		switch config.PlatformBits {
		case 64:
			asmcode += "\t;--- setup stack frame ---\n"
			asmcode += "\tpush rbp\t\t\t; save old base pointer\n"
			asmcode += "\tmov rbp, rsp\t\t\t; use stack pointer as new base pointer\n"
		case 32:
			asmcode += "\t;--- setup stack frame ---\n"
			asmcode += "\tpush ebp\t\t\t; save old base pointer\n"
			asmcode += "\tmov ebp, esp\t\t\t; use stack pointer as new base pointer\n"
		}
		return asmcode
	} else if (st[0].T == KEYWORD) && (st[0].Value == "call") && (len(st) == 2) {
		if st[1].T == VALIDNAME {
			return "\t;--- call the \"" + st[1].Value + "\" function ---\n\tcall " + st[1].Value + "\n"
		}
		log.Fatalln("Error: Calling an invalid name:", st[1].Value)
		// TODO: Find a shorter format to describe matching tokens.
		// Something along the lines of: if match(st, [KEYWORD:"extern"], 2)
	} else if (st[0].T == KEYWORD) && (st[0].Value == "counter") && (len(st) == 2) {
		return "\tmov " + config.counterRegister() + ", " + st[1].Value + "\t\t\t; set (loop) counter\n"
	} else if (st[0].T == KEYWORD) && (st[0].Value == "value") && (len(st) == 2) {
		asmcode := ""
		switch config.PlatformBits {
		case 64:
			asmcode = "\tmov rax, " + st[1].Value + "\t\t\t; set value, in preparation for looping\n"
			ps.loopStep = 8
		case 32:
			asmcode = "\tmov eax, " + st[1].Value + "\t\t\t; set value, in preparation for looping\n"
			ps.loopStep = 4
		case 16:
			// Find out if the value is a byte or a word, then set a global variable to keep track of if the nest loop should be using stosb or stosw
			if st[1].T == VALUE {
				if (strings.HasPrefix(st[1].Value, "0x") && (len(st[1].Value) == 6)) || (numbits(st[1].Value) > 8) {
					asmcode += "\tmov ax, " + st[1].Value + "\t\t\t; set value, in preparation for stosw\n"
					ps.loopStep = 2
				} else if (strings.HasPrefix(st[1].Value, "0x") && (len(st[1].Value) == 4)) || (numbits(st[1].Value) <= 8) {
					asmcode += "\tmov al, " + st[1].Value + "\t\t\t; set value, in preparation for stosb\n"
					ps.loopStep = 1
				} else {
					log.Fatalln("Error: Unable to tell if this is a word or a byte:", st[1].Value)
				}
			} else if st[1].T == REGISTER {
				switch st[1].Value {
				// TODO: Introduce a function for checking if a register is 8-bit, 16-bit, 32-bit or 64-bit
				case "al", "ah", "bl", "bh", "cl", "ch", "dl", "dh":
					asmcode += "\tmov al, " + st[1].Value + "\t\t\t; set value from register, in preparation for stosb\n"
					ps.loopStep = 1
				default:
					asmcode += "\tmov ax, " + st[1].Value + "\t\t\t; Set value from register, in preparation for stosw\n"
					ps.loopStep = 2
				}
			} else {
				log.Fatalln("Error: Unable to tell if this is a word or a byte:", st[1].Value)
			}
		default:
			log.Fatalln("Error: Unimplemented: the", st[0].Value, "keyword for", config.PlatformBits, "bit platforms")
		}
		return asmcode
	} else if (st[0].T == KEYWORD) && (st[0].Value == "loopwrite") && (len(st) == 1) {
		asmcode := ""
		switch config.PlatformBits {
		case 16:
			if ps.loopStep == 2 {
				asmcode += "\trep stosw\t\t\t; write the value in ax, cx times, starting at es:di\n"
			} else { // if ps.loop_step == 1 {
				asmcode += "\trep stosb\t\t\t; write the value in al, cx times, starting at es:di\n"
			}
		default:
			asmcode += "\tcld\n\trep stosb\t\t\t; write the value in eax/rax, ecx/rcx times, starting at edi/rdi\n"
		}
		return asmcode
	} else if (st[0].T == KEYWORD) && (st[0].Value == "write") && (len(st) == 1) {
		asmcode := ""
		switch config.PlatformBits {
		case 16:
			if ps.loopStep == 2 {
				asmcode += "\tstosw\t\t\t; write the value in ax, starting at es:di\n"
			} else { // if ps.loop_step == 1 {
				asmcode += "\tstosb\t\t\t; write the value in al, starting at es:di\n"
			}
			//else log.Fatalln("Error: Unrecognized step size. Defaulting to 1.")
		default:
			log.Fatalln("Error: Unimplemented: the", st[0].Value, "keyword for", config.PlatformBits, "bit platforms")
		}
		return asmcode
	} else if (st[0].T == KEYWORD) && ((st[0].Value == "rawloop") || (st[0].Value == "loop")) && ((len(st) == 1) || (len(st) == 2)) {
		// TODO: Make every instruction and call declare which registers they will change. This allows for better use of the registers.

		// The start of a rawloop or loop, that have an optional counter value and ends with "end"
		rawloop := (st[0].Value == "rawloop")
		hascounter := (len(st) == 2)
		endlessloop := !rawloop && !hascounter

		// Find a suitable label
		label := ""
		if rawloop {
			label = rawloopPrefix + ps.newLoopLabel()
		} else {
			if endlessloop {
				label = endlessloopPrefix + ps.newLoopLabel()
			} else {
				label = ps.newLoopLabel()
			}
		}

		// Now in the loop, in_loop is global
		ps.inLoop = label

		asmcode := ""

		// Initialize the loop, if it was given a number
		if !hascounter {
			asmcode += "\t;--- loop ---\n"
		} else {
			if endlessloop {
				asmcode += "\t;--- endless loop ---\n"
			} else {
				asmcode += "\t;--- loop " + st[1].Value + " times ---\n"
				asmcode += "\tmov " + config.counterRegister() + ", " + st[1].Value
				asmcode += "\t\t\t; initialize loop counter\n"
			}
		}
		asmcode += label + ":\t\t\t\t\t; start of loop " + label + "\n"

		// If it's not a raw loop (or endless loop), take care of the counter
		if (!rawloop) && (!endlessloop) {
			asmcode += "\tpush " + config.counterRegister() + "\t\t\t; save the counter\n"
		}
		return asmcode
	} else if (st[0].T == KEYWORD) && (st[0].Value == "address") && (len(st) == 2) {
		asmcode := ""
		switch config.PlatformBits {
		case 16:
			segmentOffset := st[1].Value
			if !strings.Contains(segmentOffset, ":") {
				log.Fatalln("Error: address takes a segment:offset value")
			}
			sl := strings.SplitN(segmentOffset, ":", 2)
			if len(sl) != 2 {
				log.Fatalln("Error: Unrecognized segment:offset address:", segmentOffset)
			}
			segment := sl[0]
			offset := sl[1]
			log.Println("Found segment", segment, "and offset", offset)
			asmcode += "\tpush " + segment + "\t\t\t; can not mov directly into es\n"
			asmcode += "\tpop es\t\t\t\t; segment = " + segment + "\n"
			// TODO: Introduce a function that checks of 0, 0x0, 0x00, 0x0000 and all other variations of zero
			if offset == "0" {
				asmcode += "\txor di, di\t\t\t; offset = " + offset + "\n"
			} else {
				asmcode += "\tmov di, " + offset + "\t\t\t; di = " + offset + "\n"
			}
		case 32:
			asmcode += "\tmov edi, " + st[1].Value + "\t\t\t; set address/offset\n"
		case 64:
			asmcode += "\tmov rdi, " + st[1].Value + "\t\t\t; set address/offset\n"
		default:
			log.Fatalln("Error: Unimplemented: the", st[0].Value, "keyword for", config.PlatformBits, "bit platforms")
		}
		return asmcode
	} else if (st[0].T == KEYWORD) && (st[0].Value == "bootable") && (len(st) == 1) {
		config.BootableKernel = true
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
		//'
	} else if (st[0].T == KEYWORD) && (st[0].Value == "extern") && (len(st) == 2) {
		if st[1].T == VALIDNAME {
			extname := st[1].Value
			// Declare the external name
			if has(ps.definedNames, extname) {
				log.Fatalln("Error: Can not declare external symbol, name is already defined: " + extname)
			}
			// Store the name of the declared constant in defined_names
			ps.definedNames = append(ps.definedNames, extname)
			// Return a comment
			return "extern " + extname + "\t\t\t; external symbol\n"
		}
		log.Fatalln("Error: extern with invalid name:", st[1].Value)
	} else if (st[0].T == KEYWORD) && (st[0].Value == "break") && (len(st) == 4) && (st[2].T == COMPARISON) {
		// breakif
		if ps.inLoop != "" {
			asmcode := ""
			rawloop := strings.HasPrefix(ps.inLoop, rawloopPrefix)     // Is it a rawloop?
			endless := strings.HasPrefix(ps.inLoop, endlessloopPrefix) // Is it endless?
			if !rawloop && !endless {
				asmcode += "\tpop " + config.counterRegister() + "\t\t\t\t; restore counter\n"
			}

			// Break if something comparison something
			asmcode += "\tcmp " + st[1].Value + ", " + st[3].Value + "\t\t\t; compare\n"

			// Conditional jump
			asmcode += "\t"
			switch st[2].Value {
			case "==":
				asmcode += "je"
			case "!=":
				asmcode += "jne"
			case ">":
				asmcode += "jg"
			case "<":
				asmcode += "jl"
			case "<=":
				asmcode += "jle"
			case ">=":
				asmcode += "jge"
			}

			// Which label to jump to (out of the loop)
			asmcode += " " + ps.inLoop + "_end\t\t\t; break\n"
			return asmcode
		}
		log.Fatalln("Error: Unclear which loop one should break out of.")
	} else if (st[0].T == KEYWORD) && (st[0].Value == "break") && (len(st) == 1) {
		if ps.inLoop != "" {
			asmcode := ""
			rawloop := strings.HasPrefix(ps.inLoop, rawloopPrefix)     // Is it a rawloop?
			endless := strings.HasPrefix(ps.inLoop, endlessloopPrefix) // Is it endless?
			if !rawloop && !endless {
				asmcode += "\tpop " + config.counterRegister() + "\t\t\t\t; restore counter\n"
			}
			asmcode += "\tjmp " + ps.inLoop + "_end\t\t\t; break\n"
			return asmcode
		}
		log.Fatalln("Error: Unclear which loop one should break out of.")
	} else if (st[0].T == KEYWORD) && (st[0].Value == "continue") && (len(st) == 4) && (st[2].T == COMPARISON) {
		// continueif
		if ps.inLoop != "" {
			asmcode := ""
			rawloop := strings.HasPrefix(ps.inLoop, rawloopPrefix)     // Is it a rawloop?
			endless := strings.HasPrefix(ps.inLoop, endlessloopPrefix) // Is it endless?
			if !rawloop && !endless {
				asmcode += "\tpop " + config.counterRegister() + "\t\t\t\t; restore counter\n"
			}

			// Continue looping if the counter is greater than zero
			//asmcode += "\tloop " + in_loop + "\t\t\t; continue\n"
			// loop can only jump <= 127 bytes away. Use dec and jnz instead
			if !endless {
				asmcode += "\tdec " + config.counterRegister() + "\t\t\t\t; decrease counter\n"
				asmcode += "\tjz " + ps.inLoop + "_end\t\t\t; jump out if the loop is done\n"
			}

			// Continue if something comparison something
			asmcode += "\tcmp " + st[1].Value + ", " + st[3].Value + "\t\t\t; compare\n"

			// Conditional jump
			asmcode += "\t"
			switch st[2].Value {
			case "==":
				asmcode += "je"
			case "!=":
				asmcode += "jne"
			case ">":
				asmcode += "jg"
			case "<":
				asmcode += "jl"
			case "<=":
				asmcode += "jle"
			case ">=":
				asmcode += "jge"
			}

			// Jump to the top if the condition is true
			asmcode += " " + ps.inLoop + "\t\t\t; continue\n"

			return asmcode
		}
		log.Fatalln("Error: Unclear which loop one should continue to the top of.")
	} else if (st[0].T == KEYWORD) && (st[0].Value == "continue") && (len(st) == 1) {
		if ps.inLoop != "" {
			asmcode := ""
			rawloop := strings.HasPrefix(ps.inLoop, rawloopPrefix)     // Is it a rawloop?
			endless := strings.HasPrefix(ps.inLoop, endlessloopPrefix) // Is it endless?
			if !rawloop && !endless {
				asmcode += "\tpop " + config.counterRegister() + "\t\t\t\t; restore counter\n"
			}
			// Continue looping if the counter is greater than zero
			//asmcode += "\tloop " + in_loop + "\t\t\t; continue\n"
			// loop can only jump <= 127 bytes away. Using dec and jnz instead
			if !endless {
				asmcode += "\tdec " + config.counterRegister() + "\t\t\t\t; decrease counter\n"
				asmcode += "\tjnz " + ps.inLoop + "\t\t\t; continue if not zero\n"
				// If the counter is zero after restoring the counter, jump out of the loop
				asmcode += "\tjz " + ps.inLoop + "_end\t\t\t; jump out if the loop is done\n"
			} else {
				asmcode += "\tjmp " + ps.inLoop + "\t\t\t; continue\n"
			}
			return asmcode
		}
		log.Fatalln("Error: Unclear which loop one should continue to the top of.")
	} else if (st[0].T == KEYWORD) && (st[0].Value == "endless") && (len(st) == 1) {
		//ps.in_loop = ""
		//ps.in_function = ""
		ps.endless = true
		return "; there is no return\n"
	} else if (st[0].T == KEYWORD) && (st[0].Value == "end") && (len(st) == 1) {
		if parseState.inlineC {
			parseState.inlineC = false
			return "; end of inline C block\n"
		} else if ps.inIfBlock != "" {
			// End the if block
			asmcode := ""
			asmcode += ps.inIfBlock + "_end:\t\t\t\t; end of if block " + ps.inIfBlock + "\n"
			ps.inIfBlock = ""
			return asmcode
		} else if ps.inLoop != "" {
			asmcode := ""
			rawloop := strings.HasPrefix(ps.inLoop, rawloopPrefix)     // Is it a rawloop?
			endless := strings.HasPrefix(ps.inLoop, endlessloopPrefix) // Is it endless?
			if !rawloop && !endless {
				asmcode += "\tpop " + config.counterRegister() + "\t\t\t\t; restore counter\n"
			}
			if endless {
				asmcode += "\tjmp " + ps.inLoop + "\t\t\t\t; loop forever\n"
				ps.endless = true
			} else {
				//asmcode += "\tloop " + in_loop + "\t\t\t\t; loop until " + config.counter_register() + " is zero\n"
				asmcode += "\tdec " + config.counterRegister() + "\t\t\t\t; decrease counter\n"
				asmcode += "\tjnz " + ps.inLoop + "\t\t\t\t; loop until " + config.counterRegister() + " is zero\n"
			}
			asmcode += ps.inLoop + "_end:\t\t\t\t; end of loop " + ps.inLoop + "\n"
			asmcode += "\t;--- end of loop " + ps.inLoop + " ---\n"
			ps.inLoop = ""
			return asmcode
		} else if ps.inFunction != "" {
			// Return from the function if "end" is encountered
			ret := Token{KEYWORD, "ret", st[0].Line, ""}
			newstatement := Statement{ret}
			return newstatement.String(ps, config)
		} else {
			// If the function was already ended with "exit", don't freak out when encountering an "end"
			if !ps.surpriseEndingWithExit && !ps.endless {
				log.Fatalln("Error: Not in a function or block of inline C, hard to tell what should be ended with \"end\". Statement nr:", st[0].Line)
			} else {
				// Prepare for more surprises
				ps.surpriseEndingWithExit = false
				// Ignore this "end"
				return ""
			}
		}
	} else if (st[0].T == VALIDNAME) && (len(st) == 1) {
		// Just a name, assume it's a function call
		if has(ps.definedNames, st[0].Value) {
			call := Token{KEYWORD, "call", st[0].Line, ""}
			newstatement := Statement{call, st[0]}
			return newstatement.String(ps, config)
		}
		log.Fatalln("Error: No function named:", st[0].Value)
	} else if (st[0].T == KEYWORD) && (st[0].Value == "noret") {
		return "; end without a return\n"
	} else if (st[0].T == KEYWORD) && (st[0].Value == "inline_c") {
		parseState.inlineC = true
		return "; start of inline C block\n"
	} else if (st[0].T == KEYWORD) && (st[0].Value == "const") {
		log.Fatalln("Error: Incomprehensible constant:", st.String(ps, config))
	} else if st[0].T == BUILTIN {
		log.Fatalln("Error: Unhandled builtin:", st[0].Value)
	} else if st[0].T == KEYWORD {
		log.Fatalln("Error: Unhandled keyword:", st[0].Value)
	}
	log.Println("Error: Unfamiliar statement layout: ")
	for _, token := range []Token(st) {
		log.Print(token)
	}
	os.Exit(1)
	return ";ERROR"
}
