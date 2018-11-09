package ez

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Type int

const (
	TYP_INT Type = iota
	TYP_STR
	TYP_ANY
)

type Prog struct {
	Pos         int
	OpAddrs     []int
	Ints        []int
	Strs        []string
	IntIDAddr   map[string]int
	StrIDAddr   map[string]int
	FuncIDImpls map[string][]Func
	FuncAddrs   []func(*Prog)
}

type Func struct {
	Variadic bool
	In       []Type
	Out      []Type
	F        func(*Prog)
	addr     int
}

func Prelude() map[string][]Func {
	prelude := map[string][]Func{}
	prelude["+"] = []Func{
		{
			Variadic: false,
			In:       []Type{TYP_INT, TYP_INT},
			Out:      []Type{TYP_INT},
			F: func(p *Prog) {
				p.Ints[p.OpAddrs[p.Pos+3]] = p.Ints[p.OpAddrs[p.Pos+1]] + p.Ints[p.OpAddrs[p.Pos+2]]
				p.Pos += 3
			},
		},
		{
			Variadic: false,
			In:       []Type{TYP_STR, TYP_STR},
			Out:      []Type{TYP_STR},
			F: func(p *Prog) {
				p.Strs[p.OpAddrs[p.Pos+3]] = p.Strs[p.OpAddrs[p.Pos+1]] + p.Strs[p.OpAddrs[p.Pos+2]]
				p.Pos += 3
			},
		},
	}
	prelude["*"] = []Func{
		{
			Variadic: false,
			In:       []Type{TYP_INT, TYP_INT},
			Out:      []Type{TYP_INT},
			F: func(p *Prog) {
				p.Ints[p.OpAddrs[p.Pos+3]] = p.Ints[p.OpAddrs[p.Pos+1]] * p.Ints[p.OpAddrs[p.Pos+2]]
				p.Pos += 3
			},
		},
	}
	prelude["/"] = []Func{
		{
			Variadic: false,
			In:       []Type{TYP_INT, TYP_INT},
			Out:      []Type{TYP_INT},
			F: func(p *Prog) {
				p.Ints[p.OpAddrs[p.Pos+3]] = p.Ints[p.OpAddrs[p.Pos+1]] / p.Ints[p.OpAddrs[p.Pos+2]]
				p.Pos += 3
			},
		},
	}
	prelude["-"] = []Func{
		{
			Variadic: false,
			In:       []Type{TYP_INT, TYP_INT},
			Out:      []Type{TYP_INT},
			F: func(p *Prog) {
				p.Ints[p.OpAddrs[p.Pos+3]] = p.Ints[p.OpAddrs[p.Pos+1]] - p.Ints[p.OpAddrs[p.Pos+2]]
				p.Pos += 3
			},
		},
	}
	return prelude
}

const (
	IOP_INTCOPY = iota
	IOP_STRCOPY
)

var internalOps = []func(*Prog){
	func(p *Prog) {
		p.Ints[p.OpAddrs[p.Pos+2]] = p.Ints[p.OpAddrs[p.Pos+1]]
		p.Pos += 2
	},
	func(p *Prog) {
		p.Strs[p.OpAddrs[p.Pos+2]] = p.Strs[p.OpAddrs[p.Pos+1]]
		p.Pos += 2
	},
}

func (p *Prog) Run() {
	for p.Pos < len(p.OpAddrs) {
		p.FuncAddrs[p.OpAddrs[p.Pos]](p)
		p.Pos += 1
	}
	p.Pos = 0
}

func (p *Prog) AddOps(ops map[string][]Func) error {
	if len(p.FuncAddrs) == 0 {
		p.FuncAddrs = internalOps
	}
	for opName, opFuncs := range ops {
		for i, opFunc := range opFuncs {
			opFunc.addr = len(p.FuncAddrs)
			p.FuncAddrs = append(p.FuncAddrs, opFunc.F)
			opFuncs[i] = opFunc
		}
		if existingFuncs, ok := p.FuncIDImpls[opName]; ok {
			p.FuncIDImpls[opName] = append(existingFuncs, opFuncs...)
		} else {
			p.FuncIDImpls[opName] = opFuncs
		}
	}
	return nil
}

func (p *Prog) DebugPrintSymbol(symbol string) {
	typ, addr, ok := p.typeAndAddrOfID(symbol)
	if !ok {
		println("ERR - p.DebugPrintSymbol: no symbol with name '" + symbol + "'")
	}
	switch typ {
	case TYP_STR:
		println(p.Strs[addr])
	case TYP_INT:
		println(p.Ints[addr])
	}
}

func (p *Prog) Parse(input []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	scanner.Split(bufio.ScanLines)

	var line int

	for scanner.Scan() {
		line += 1

		var assgns []string
		var buildingAssgns bool
		var args []string
		var op string

		var buildStr string
		var buildingStr bool

		fields := strings.Fields(scanner.Text())
		for i, field := range fields {
			if strings.HasPrefix(field, "#") {
				break
			}
			switch {
			case buildingStr:
				buildStr += " " + field
				if isStringEnd(field) {
					buildingStr = false
					args = append(args, buildStr)
				}
			case i == 0 && isIdentifier(field):
				assgns = append(assgns, field)
				buildingAssgns = true
			case field == "=":
				if !buildingAssgns {
					return parsingErr(line, "expected one or more identifiers to left of assigment operator")
				}
				buildingAssgns = false
			case p.isFuncCall(field):
				op = field
			case isStringStart(field):
				if len(field) >= 2 && isStringEnd(field) {
					args = append(args, field)
				} else {
					buildStr = field
					buildingStr = true
				}
			case isIdentifier(field) || isInt(field):
				if buildingAssgns {
					if !isIdentifier(field) {
						return parsingErr(line, "expected another identifier or an assignment symbol '=', got '"+field+"'")
					}
					assgns = append(assgns, field)
				} else {
					args = append(args, field)
				}
			default:
				return parsingErr(line, "unknown symbol: "+field)
			}
		}
		switch {
		case op == "" && len(args) > 0 && len(assgns) > 0:
			if len(args) > 1 || len(assgns) > 1 {
				return parsingErr(line, "can only assign one expression to one argument")
			}
			targetTyp, targetAddr, targetFound := p.typeAndAddrOfID(assgns[0])
			if isIdentifier(args[0]) {
				typ, addr, found := p.typeAndAddrOfID(args[0])
				if !found {
					return parsingErr(line, "reference to uninitialized identifier: "+args[0])
				}
				if targetFound {
					if targetTyp != typ {
						return parsingErr(line, "cannot assign '"+args[0]+"' to '"+assgns[0]+"' - type mismatch")
					}
				} else {
					targetAddr = p.newAlloc(assgns[0], typ)
				}
				p.OpAddrs = append(p.OpAddrs, p.copyFuncInstructionForType(typ), addr, targetAddr)
			} else {
				if targetFound {
					if targetTyp != rawToType(args[0]) {
						return parsingErr(line, "cannot assign '"+args[0]+"' to '"+assgns[0]+"' - type mismatch")
					}
					_, addr, found := p.typeAndAddrOfID(args[0])
					if !found {
						addr = p.newAllocInitialize(args[0], args[0])
					}
					p.OpAddrs = append(p.OpAddrs, p.copyFuncInstructionForType(targetTyp), addr, targetAddr)
				} else {
					p.newAllocInitialize(assgns[0], args[0])
				}
			}
		case op != "":
			funcs, ok := p.FuncIDImpls[op]
			if !ok {
				panic("shouldn't be able to get here as op is already checked for existence during parsing: " + op)
			}
			var argTypes []Type
			var argAddrs []int
			for _, arg := range args {
				if isIdentifier(arg) {
					typ, addr, found := p.typeAndAddrOfID(arg)
					if !found {
						return parsingErr(line, "reference to uninitialized identifier: "+arg)
					}
					argTypes = append(argTypes, typ)
					argAddrs = append(argAddrs, addr)
				} else {
					typ, addr, found := p.typeAndAddrOfID(arg)
					if !found {
						addr = p.newAllocInitialize(arg, arg)
						typ = rawToType(arg) // TODO-OPT: return type from newAllocInitialize
					}
					argTypes = append(argTypes, typ)
					argAddrs = append(argAddrs, addr)
				}
			}
			var assgnTypes []Type
			var assgnAddrs []int
			for _, assgn := range assgns {
				typ, addr, found := p.typeAndAddrOfID(assgn)
				if !found {
					typ = TYP_ANY
					addr = -1
				}
				assgnTypes = append(assgnTypes, typ)
				assgnAddrs = append(assgnAddrs, addr)
			}
			foundFunc := false
			for _, fun := range funcs {
				if len(fun.In) != len(args) || len(fun.Out) != len(assgns) {
					continue
				}
				inTypesMatch := true
				for i, inType := range fun.In {
					if inType == TYP_ANY {
						continue
					}
					if inType != argTypes[i] {
						inTypesMatch = false
						break
					}
				}
				if !inTypesMatch {
					continue
				}
				outTypesMatch := true
				for i, outType := range fun.Out {
					if assgnTypes[i] == TYP_ANY {
						continue
					}
					if outType != assgnTypes[i] {
						outTypesMatch = false
						break
					}
				}
				if !outTypesMatch {
					continue
				}
				for i, outType := range fun.Out {
					if assgnTypes[i] != TYP_ANY {
						continue
					}
					assgnAddrs[i] = p.newAlloc(assgns[i], outType)
				}
				foundFunc = true
				p.OpAddrs = append(p.OpAddrs, fun.addr)
				p.OpAddrs = append(p.OpAddrs, argAddrs...)
				p.OpAddrs = append(p.OpAddrs, assgnAddrs...)
				break
			}
			if !foundFunc {
				return parsingErr(line, "no function signature named '"+op+"' to handle types/quantity of arguments or assignments")
			}
		}
	}
	return nil
}

func rawToType(raw string) Type {
	switch {
	case isString(raw):
		return TYP_STR
	case isInt(raw):
		return TYP_INT
	}
	panic("no type for: " + raw)
}

func (p *Prog) newAllocInitialize(id, raw string) int {
	var addr int
	switch rawToType(raw) {
	case TYP_STR:
		addr = len(p.Strs)
		p.StrIDAddr[id] = len(p.Strs)
		p.Strs = append(p.Strs, raw[1:len(raw)-1])
	case TYP_INT:
		convInt, err := strconv.Atoi(raw)
		if err != nil {
			panic("failed to convert int '" + raw + "' even though it was parsed as an int: " + err.Error())
		}
		addr = len(p.Ints)
		p.IntIDAddr[id] = len(p.Ints)
		p.Ints = append(p.Ints, convInt)
	}
	return addr
}

func isString(raw string) bool {
	return len(raw) > 1 && raw[0] == '\'' && raw[len(raw)-1] == '\''
}

func (p *Prog) copyFuncInstructionForType(typ Type) int {
	switch typ {
	case TYP_INT:
		return IOP_INTCOPY
	case TYP_STR:
		return IOP_STRCOPY
	}
	panic("type has no copy instruction: " + strconv.Itoa(int(typ)))
}

func (p *Prog) newAlloc(id string, typ Type) int {
	var addr int
	switch typ {
	case TYP_INT:
		addr = len(p.Ints)
		p.IntIDAddr[id] = len(p.Ints)
		p.Ints = append(p.Ints, 0)
	case TYP_STR:
		addr = len(p.Strs)
		p.StrIDAddr[id] = len(p.Strs)
		p.Strs = append(p.Strs, "")
	}
	return addr
}

func (p *Prog) copyToExisting(fromAddr, toAddr int, typ Type) {
	switch typ {
	case TYP_STR:
		p.Strs[toAddr] = p.Strs[fromAddr]
	case TYP_INT:
		p.Ints[toAddr] = p.Ints[fromAddr]
	}
}

func parsingErr(line int, errMsg string) error {
	return errors.New("ERROR - Line " + strconv.Itoa(line) + ": " + errMsg)
}

func (p *Prog) typeAndAddrOfID(id string) (Type, int, bool) {
	for typ, idAddr := range map[Type]map[string]int{
		TYP_STR: p.StrIDAddr,
		TYP_INT: p.IntIDAddr,
	} {
		if addr, ok := idAddr[id]; ok {
			return typ, addr, ok
		}
	}
	return TYP_ANY, -1, false
}

func isStringStart(str string) bool {
	for i, r := range str {
		if i == 0 {
			if r != '\'' {
				return false
			}
		}
		return true
	}
	return false
}

func isStringEnd(str string) bool {
	return len(str) > 0 && str[len(str)-1] == '\'' && (len(str) == 1 || str[len(str)-2] != '\\')
}

func isInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func isIdentifier(str string) bool {
	for i, r := range str {
		if i == 0 {
			if !unicode.IsLetter(r) || !unicode.IsLower(r) {
				return false
			}
		} else if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func (p *Prog) isFuncCall(str string) bool {
	_, ok := p.FuncIDImpls[str]
	return ok
}
