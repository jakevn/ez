package ez

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var (
	MaxLineLen        = 120
	MaxLines   uint16 = 800
)

type baseType int

const (
	Und baseType = iota
	Int
	Str
	Bool
	Any
	Addr
	ArrUnd
	ArrInt
	ArrStr
	ArrBool
)

type Parser struct {
	bc                  Bytecode
	IDInfo              map[string]Info
	UndecidedDependents map[string][]string
	InParams            map[string]Param
	OutParams           map[string]Param
	undecidedAddrIndex  int
	line                uint16
}

type Info struct {
	Type      baseType
	Addresses []Address
}

type Address struct {
	Index int
	Line  uint16
}

type Param struct {
	Pos  int
	Type baseType
	Addr int
}

type expressionCtx struct {
	assgns []string
	args   []string
	op     string
	array  bool
}

func Parse(reader io.Reader) (Bytecode, error) {
	p := &Parser{
		IDInfo:              map[string]Info{},
		InParams:            map[string]Param{},
		UndecidedDependents: map[string][]string{},
		undecidedAddrIndex:  -100,
	}
	return p.parseInternal(reader)
}

func (p *Parser) parseInternal(reader io.Reader) (Bytecode, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		p.line += 1
		if p.line > MaxLines {
			return p.bc, errors.New("exceeded max number of lines: " + strconv.Itoa(int(MaxLines)))
		}
		lineText := scanner.Text()
		if len(lineText) == 0 || lineText[0] == '#' {
			continue
		}
		if len(lineText) > MaxLineLen {
			return p.bc, p.parsingErr("exceeded max line length: " + strconv.Itoa(MaxLineLen))
		}

		var baseExprCtx expressionCtx
		var buildStr string
		var buildingStr bool
		var buildingAssgns bool
		var injectEndAddrAt int
		fields := strings.Fields(lineText)
		for i, field := range fields {
			if !buildingStr && strings.HasPrefix(field, "#") {
				break
			}
			if i == 0 && isLabel(field) {
				p.newAlloc(field, Addr)
				if len(fields) > 1 && !strings.HasPrefix(fields[1], "#") {
					return p.bc, p.parsingErr("labels can only be followed by a comment")
				}
				break
			}

			if strings.HasPrefix(field, "[") {
				if baseExprCtx.array == true {
					return p.bc, p.parsingErr("errant '[' in array declaration")
				}
				baseExprCtx.array = true
				field = field[1:]
				if len(field) == 1 {
					continue
				}
			}
			if baseExprCtx.array && strings.HasSuffix(field, "]") {
				field = field[:len(field)-2]
			}

			switch {
			case buildingStr:
				buildStr += " " + field
				if isStringEnd(field) {
					buildingStr = false
					baseExprCtx.args = append(baseExprCtx.args, buildStr)
				}
			case i == 0 && isIdentifier(field):
				baseExprCtx.assgns = append(baseExprCtx.assgns, field)
				buildingAssgns = true
			case field == "=":
				if !buildingAssgns {
					return p.bc, p.parsingErr("expected one or more identifiers to left of assigment operator")
				}
				buildingAssgns = false
			case isFuncCall(field):
				baseExprCtx.op = field
			case isStringStart(field):
				if len(field) >= 2 && isStringEnd(field) {
					baseExprCtx.args = append(baseExprCtx.args, field)
				} else {
					buildStr = field
					buildingStr = true
				}
			case isIdentifier(field) || isBool(field) || isInt(field) || isLabel(field):
				if buildingAssgns {
					if !isIdentifier(field) && !isLabel(field) {
						return p.bc, p.parsingErr("expected another identifier or an assignment symbol '=', got '" + field + "'")
					}
					baseExprCtx.assgns = append(baseExprCtx.assgns, field)
				} else {
					baseExprCtx.args = append(baseExprCtx.args, field)
				}
				if baseExprCtx.op == "if" {
					if err := p.compileExpression(baseExprCtx); err != nil {
						return p.bc, err
					}
					injectEndAddrAt = len(p.bc.OpAddrs)
					p.bc.OpAddrs = append(p.bc.OpAddrs, 0)
					baseExprCtx = expressionCtx{}
					buildingAssgns = false
				}
			default:
				return p.bc, p.parsingErr("unknown symbol: " + field)
			}
		}
		if err := p.compileExpression(baseExprCtx); err != nil {
			return p.bc, err
		}
		if injectEndAddrAt != 0 {
			p.bc.OpAddrs[injectEndAddrAt] = len(p.bc.OpAddrs)
			injectEndAddrAt = 0
		}
	}
	return p.bc, nil
}

func (p *Parser) compileExpression(ctx expressionCtx) error {
	switch {
	case ctx.array: // Array
		//if len(ctx.assgns) != 1 {
		//	return p.parsingErr("array must be assigned to a single identifier")
		//}
		//targetTyp, targetAddr, targetFound := p.typeAndAddrOfID(ctx.assgns[0])
		//prevTyp := -1
		//firstIsUnd := false
		//for i, arg := range ctx.args {
		//	if isIdentifier(arg) {
		//		typ, addr, found := p.typeAndAddrOfID(arg)
		//		if !found {
		//			return p.parsingErr("reference to uninitialized identifier: " + arg)
		//		}
		//	}
		//	p.typeAndAddrOfID(arg)
		//	if i == 0 {
		//		continue
		//	}
		//}

	case ctx.op == "" && len(ctx.args) > 0 && len(ctx.assgns) > 0:
		if len(ctx.args) > 1 || len(ctx.assgns) > 1 {
			return p.parsingErr("can only assign one expression to one argument")
		}
		targetTyp, targetAddr, targetFound := p.typeAndAddrOfID(ctx.assgns[0])
		if isIdentifier(ctx.args[0]) {
			typ, addr, found := p.typeAndAddrOfID(ctx.args[0])
			if !found {
				return p.parsingErr("reference to uninitialized identifier: " + ctx.args[0])
			}
			if targetFound {
				if targetTyp != typ {
					if typ == Und {
						targetAddr = p.undecidedIsDecided(ctx.args[0], targetTyp)
					} else {
						return p.parsingErr("cannot assign '" + ctx.args[0] + "' to '" + ctx.assgns[0] + "' - type mismatch")
					}
				}
			} else {
				if typ == Und {
					p.undecidedAddDependency(ctx.args[0], ctx.assgns[0])
				}
				targetAddr = p.newAlloc(ctx.assgns[0], typ)
			}
			p.bc.OpAddrs = append(p.bc.OpAddrs, p.copyFuncInstructionForType(typ), addr, targetAddr)
		} else {
			if targetFound {
				if targetTyp != rawToType(ctx.args[0]) {
					return p.parsingErr("cannot assign '" + ctx.args[0] + "' to '" + ctx.assgns[0] + "' - type mismatch")
				}
				_, addr, found := p.typeAndAddrOfID(ctx.args[0])
				if !found {
					addr, _ = p.newAllocInitialize(ctx.args[0], ctx.args[0])
				}
				p.bc.OpAddrs = append(p.bc.OpAddrs, p.copyFuncInstructionForType(targetTyp), addr, targetAddr)
			} else {
				p.newAllocInitialize(ctx.assgns[0], ctx.args[0])
			}
		}
	case len(ctx.assgns) > 0 && len(ctx.args) == 0 && ctx.op == "":
		if len(p.InParams) > 0 {
			return p.parsingErr("expected assignment or expression following identifier")
		}
		for i, inParamID := range ctx.assgns {
			if _, ok := p.InParams[inParamID]; ok {
				return p.parsingErr("in parameter identifiers must be unique - duplicate: '" + inParamID + "'")
			}
			p.InParams[inParamID] = Param{
				Pos: i,
			}
		}
	case ctx.op != "":
		funcs, ok := baselib[ctx.op]
		if !ok {
			return p.parsingErr("impossible made possible - previously existing op no longer exists: " + ctx.op)
		}
		var argTypes []baseType
		var argAddrs []int
		for _, arg := range ctx.args {
			if isIdentifier(arg) || isLabel(arg) {
				typ, addr, found := p.typeAndAddrOfID(arg)
				if !found {
					return p.parsingErr("reference to uninitialized identifier: " + arg)
				}
				argTypes = append(argTypes, typ)
				argAddrs = append(argAddrs, addr)
			} else {
				addr, typ := p.newAllocInitialize(arg, arg)
				argTypes = append(argTypes, typ)
				argAddrs = append(argAddrs, addr)
			}
		}
		var assgnTypes []baseType
		var assgnAddrs []int
		for _, assgn := range ctx.assgns {
			typ, addr, found := p.typeAndAddrOfID(assgn)
			if !found {
				typ = Any
				addr = -1
			}
			assgnTypes = append(assgnTypes, typ)
			assgnAddrs = append(assgnAddrs, addr)
		}
		foundFunc := false
		for _, fun := range funcs {
			if len(fun.In) != len(ctx.args) || len(fun.Out) != len(ctx.assgns) {
				continue // TODO: overlapping func names can no longer have diff param/return len
			}
			inTypesMatch := true
			for i, inType := range fun.In {
				if inType == Und {
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
				if assgnTypes[i] == Und || assgnTypes[i] == Any {
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
				switch assgnTypes[i] {
				case Any:
					assgnAddrs[i] = p.newAlloc(ctx.assgns[i], outType)
				case Und:
					assgnAddrs[i] = p.undecidedIsDecided(ctx.assgns[i], outType)
				}
			}
			for i, inType := range fun.In {
				if argTypes[i] != Und {
					continue
				}
				argAddrs[i] = p.undecidedIsDecided(ctx.args[i], inType)
			}
			foundFunc = true
			p.bc.OpAddrs = append(p.bc.OpAddrs, fun.addr)
			p.bc.OpAddrs = append(p.bc.OpAddrs, argAddrs...)
			p.bc.OpAddrs = append(p.bc.OpAddrs, assgnAddrs...)
			break
		}
		if !foundFunc {
			return p.parsingErr("no function signature named '" + ctx.op + "' to handle types/quantity of arguments or assignments")
		}
	}
	return nil
}

func (p *Parser) newOrGetUndecidedAddr(id string) int {
	undecided, ok := p.IDInfo[id]
	if !ok {
		p.undecidedAddrIndex -= 1
		p.IDInfo[id] = Info{
			Type: Und,
			Addresses: []Address{
				{Line: p.line, Index: p.undecidedAddrIndex},
			},
		}
		return p.undecidedAddrIndex
	}
	return undecided.Addresses[len(undecided.Addresses)-1].Index
}

func (p *Parser) undecidedAddDependency(id, dependentId string) {
	dependents, ok := p.UndecidedDependents[id]
	if !ok {
		p.UndecidedDependents[id] = []string{dependentId}
		return
	}
	dependents = append(dependents, dependentId)
	p.UndecidedDependents[id] = dependents
}

func (p *Parser) undecidedIsDecided(id string, typ baseType) int {
	undecided, ok := p.IDInfo[id]
	if !ok || undecided.Type != Und {
		return -1
	}
	undecided.Type = typ
	p.IDInfo[id] = undecided
	if dependents, ok := p.UndecidedDependents[id]; ok {
		for _, dep := range dependents {
			p.undecidedIsDecided(dep, typ)
		}
	}
	var latestAddr int
	for i, addr := range undecided.Addresses {
		latestAddr = p.newAlloc(id, typ)
		if inputParam, ok := p.InParams[id]; ok && inputParam.Type == Und {
			inputParam.Type = typ
			inputParam.Addr = latestAddr
			p.InParams[id] = inputParam
		}
		for j, opAddr := range p.bc.OpAddrs {
			if addr.Index == opAddr {
				p.bc.OpAddrs[j] = latestAddr
			}
		}
		undecided.Addresses[i].Index = latestAddr
	}
	return latestAddr
}

func (p *Parser) newAllocInitialize(id, raw string) (int, baseType) {
	var addr int
	typ := rawToType(raw)
	switch typ {
	case Str:
		addr = len(p.bc.Strs)
		p.bc.Strs = append(p.bc.Strs, raw[1:len(raw)-1])
	case Int:
		convInt, err := strconv.Atoi(raw)
		if err != nil {
			panic("failed to convert int '" + raw + "' even though it was parsed as an int: " + err.Error())
		}
		addr = len(p.bc.Ints)
		p.bc.Ints = append(p.bc.Ints, convInt)
	case Bool:
		addr = len(p.bc.Bools)
		p.bc.Bools = append(p.bc.Bools, raw == "true")
	}
	p.IDInfo[id] = Info{
		Type:      typ,
		Addresses: []Address{{Index: addr, Line: p.line}},
	}
	return addr, typ
}

// TODO: handle copy instructions for undecided type
func (p *Parser) copyFuncInstructionForType(typ baseType) int {
	switch typ {
	case Int:
		return iopIntCopy
	case Str:
		return iopStrCopy
	case Bool:
		return iopBoolCopy
	}
	panic("type has no copy instruction: " + strconv.Itoa(int(typ)))
}

func (p *Parser) newAlloc(id string, typ baseType) int {
	var addr int
	switch typ {
	case Int:
		addr = len(p.bc.Ints)
		p.bc.Ints = append(p.bc.Ints, 0)
	case Str:
		addr = len(p.bc.Strs)
		p.bc.Strs = append(p.bc.Strs, "")
	case Bool:
		addr = len(p.bc.Bools)
		p.bc.Bools = append(p.bc.Bools, false)
	case Und:
		addr = p.newOrGetUndecidedAddr(id) // TODO
	case Addr:
		addr = len(p.bc.Ints)
		p.bc.Ints = append(p.bc.Ints, len(p.bc.OpAddrs))
	}
	p.IDInfo[id] = Info{
		Type:      typ,
		Addresses: []Address{{Index: addr, Line: p.line}},
	}
	return addr
}

func (p *Parser) parsingErr(errMsg string) error {
	return errors.New("ERROR - line " + strconv.Itoa(int(p.line)) + ": " + errMsg)
}

func (p *Parser) typeAndAddrOfID(id string) (baseType, int, bool) {
	if info, ok := p.IDInfo[id]; ok {
		return info.Type, info.Addresses[len(info.Addresses)-1].Index, ok
	}
	return Und, -1, false
}

func rawToType(raw string) baseType {
	switch {
	case isString(raw):
		return Str
	case isInt(raw):
		return Int
	case isBool(raw):
		return Bool
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
	if isFuncCall(str) {
		return false
	}
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

func isLabel(str string) bool {
	return str[0] == '~' && len(str) > 1 && isIdentifier(str[1:])
}

func isString(raw string) bool {
	return len(raw) > 1 && raw[0] == '\'' && raw[len(raw)-1] == '\''
}

func isBool(str string) bool {
	return str == "true" || str == "false"
}

func isFuncCall(str string) bool {
	_, ok := baselib[str]
	return ok
}
