package ez

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type baseType int

const (
	btInt baseType = iota
	btStr
	btBool
	btAny
)

type Parser struct {
	bc         Bytecode
	IntIDAddr  map[string]int
	StrIDAddr  map[string]int
	BoolIDAddr map[string]int
}

type Func struct {
	In   []baseType
	Out  []baseType
	F    func(*Bytecode)
	addr int
}

func Parse(reader io.Reader) (Bytecode, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	p := &Parser{
		IntIDAddr:  map[string]int{},
		StrIDAddr:  map[string]int{},
		BoolIDAddr: map[string]int{},
	}
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
					return p.bc, parsingErr(line, "expected one or more identifiers to left of assigment operator")
				}
				buildingAssgns = false
			case isFuncCall(field):
				op = field
			case isStringStart(field):
				if len(field) >= 2 && isStringEnd(field) {
					args = append(args, field)
				} else {
					buildStr = field
					buildingStr = true
				}
			case isIdentifier(field) || isBool(field) || isInt(field):
				if buildingAssgns {
					if !isIdentifier(field) {
						return p.bc, parsingErr(line, "expected another identifier or an assignment symbol '=', got '"+field+"'")
					}
					assgns = append(assgns, field)
				} else {
					args = append(args, field)
				}
			default:
				return p.bc, parsingErr(line, "unknown symbol: "+field)
			}
		}
		switch {
		case op == "" && len(args) > 0 && len(assgns) > 0:
			if len(args) > 1 || len(assgns) > 1 {
				return p.bc, parsingErr(line, "can only assign one expression to one argument")
			}
			targetTyp, targetAddr, targetFound := p.typeAndAddrOfID(assgns[0])
			if isIdentifier(args[0]) {
				typ, addr, found := p.typeAndAddrOfID(args[0])
				if !found {
					return p.bc, parsingErr(line, "reference to uninitialized identifier: "+args[0])
				}
				if targetFound {
					if targetTyp != typ {
						return p.bc, parsingErr(line, "cannot assign '"+args[0]+"' to '"+assgns[0]+"' - type mismatch")
					}
				} else {
					targetAddr = p.newAlloc(assgns[0], typ)
				}
				p.bc.OpAddrs = append(p.bc.OpAddrs, p.copyFuncInstructionForType(typ), addr, targetAddr)
			} else {
				if targetFound {
					if targetTyp != rawToType(args[0]) {
						return p.bc, parsingErr(line, "cannot assign '"+args[0]+"' to '"+assgns[0]+"' - type mismatch")
					}
					_, addr, found := p.typeAndAddrOfID(args[0])
					if !found {
						addr, _ = p.newAllocInitialize(args[0], args[0])
					}
					p.bc.OpAddrs = append(p.bc.OpAddrs, p.copyFuncInstructionForType(targetTyp), addr, targetAddr)
				} else {
					p.newAllocInitialize(assgns[0], args[0])
				}
			}
		case op != "":
			funcs, ok := baselib[op]
			if !ok {
				panic("shouldn't be able to get here as op is already checked for existence during parsing: " + op)
			}
			var argTypes []baseType
			var argAddrs []int
			for _, arg := range args {
				if isIdentifier(arg) {
					typ, addr, found := p.typeAndAddrOfID(arg)
					if !found {
						return p.bc, parsingErr(line, "reference to uninitialized identifier: "+arg)
					}
					argTypes = append(argTypes, typ)
					argAddrs = append(argAddrs, addr)
				} else {
					typ, addr, found := p.typeAndAddrOfID(arg)
					if !found {
						addr, typ = p.newAllocInitialize(arg, arg)
					}
					argTypes = append(argTypes, typ)
					argAddrs = append(argAddrs, addr)
				}
			}
			var assgnTypes []baseType
			var assgnAddrs []int
			for _, assgn := range assgns {
				typ, addr, found := p.typeAndAddrOfID(assgn)
				if !found {
					typ = btAny
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
					if inType == btAny {
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
					if assgnTypes[i] == btAny {
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
					if assgnTypes[i] != btAny {
						continue
					}
					assgnAddrs[i] = p.newAlloc(assgns[i], outType)
				}
				foundFunc = true
				p.bc.OpAddrs = append(p.bc.OpAddrs, fun.addr)
				p.bc.OpAddrs = append(p.bc.OpAddrs, argAddrs...)
				p.bc.OpAddrs = append(p.bc.OpAddrs, assgnAddrs...)
				break
			}
			if !foundFunc {
				return p.bc, parsingErr(line, "no function signature named '"+op+"' to handle types/quantity of arguments or assignments")
			}
		}
	}
	return p.bc, nil
}

func (p *Parser) newAllocInitialize(id, raw string) (int, baseType) {
	var addr int
	var typ baseType
	switch rawToType(raw) {
	case btStr:
		typ = btStr
		addr = len(p.bc.Strs)
		p.StrIDAddr[id] = len(p.bc.Strs)
		p.bc.Strs = append(p.bc.Strs, raw[1:len(raw)-1])
	case btInt:
		typ = btInt
		convInt, err := strconv.Atoi(raw)
		if err != nil {
			panic("failed to convert int '" + raw + "' even though it was parsed as an int: " + err.Error())
		}
		addr = len(p.bc.Ints)
		p.IntIDAddr[id] = len(p.bc.Ints)
		p.bc.Ints = append(p.bc.Ints, convInt)
	case btBool:
		typ = btBool
		addr = len(p.bc.Bools)
		p.BoolIDAddr[id] = len(p.bc.Bools)
		p.bc.Bools = append(p.bc.Bools, raw == "True")
	}
	return addr, typ
}

func (p *Parser) copyFuncInstructionForType(typ baseType) int {
	switch typ {
	case btInt:
		return iopIntCopy
	case btStr:
		return iopStrCopy
	case btBool:
		return iopBoolCopy
	}
	panic("type has no copy instruction: " + strconv.Itoa(int(typ)))
}

func (p *Parser) newAlloc(id string, typ baseType) int {
	var addr int
	switch typ {
	case btInt:
		addr = len(p.bc.Ints)
		p.IntIDAddr[id] = len(p.bc.Ints)
		p.bc.Ints = append(p.bc.Ints, 0)
	case btStr:
		addr = len(p.bc.Strs)
		p.StrIDAddr[id] = len(p.bc.Strs)
		p.bc.Strs = append(p.bc.Strs, "")
	case btBool:
		addr = len(p.bc.Bools)
		p.BoolIDAddr[id] = len(p.bc.Bools)
		p.bc.Bools = append(p.bc.Bools, false)
	}
	return addr
}

func (p *Parser) copyToExisting(fromAddr, toAddr int, typ baseType) {
	switch typ {
	case btStr:
		p.bc.Strs[toAddr] = p.bc.Strs[fromAddr]
	case btInt:
		p.bc.Ints[toAddr] = p.bc.Ints[fromAddr]
	case btBool:
		p.bc.Bools[toAddr] = p.bc.Bools[fromAddr]
	}
}

func parsingErr(line int, errMsg string) error {
	return errors.New("ERROR - Line " + strconv.Itoa(line) + ": " + errMsg)
}

func (p *Parser) typeAndAddrOfID(id string) (baseType, int, bool) {
	for typ, idAddr := range map[baseType]map[string]int{
		btStr:  p.StrIDAddr,
		btInt:  p.IntIDAddr,
		btBool: p.BoolIDAddr,
	} {
		if addr, ok := idAddr[id]; ok {
			return typ, addr, ok
		}
	}
	return btAny, -1, false
}

func rawToType(raw string) baseType {
	switch {
	case isString(raw):
		return btStr
	case isInt(raw):
		return btInt
	case isBool(raw):
		return btBool
	}
	panic("no type for: " + raw)
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

func isString(raw string) bool {
	return len(raw) > 1 && raw[0] == '\'' && raw[len(raw)-1] == '\''
}

func isBool(str string) bool {
	return str == "True" || str == "False"
}

func isFuncCall(str string) bool {
	_, ok := baselib[str]
	return ok
}
