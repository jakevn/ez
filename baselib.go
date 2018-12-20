package ez

import (
	"log"
	"sort"
)

func init() {
	var syms []string
	for sym := range baselib {
		syms = append(syms, sym)
	}
	sort.Strings(syms)
	for _, sym := range syms {
		symFuncs := baselib[sym]
		if len(symFuncs) == 0 {
			panic("symbol with no associated functions: " + sym)
		}
		paramInCount := len(symFuncs[0].In)
		paramOutCount := len(symFuncs[0].Out)
		for i, fun := range symFuncs {
			if len(fun.In) != paramInCount || len(fun.Out) != paramOutCount {
				panic("functions associated with symbols must have identical parameter counts. invalid: " + sym)
			}
			fun.addr = len(funcAddrs)
			funcAddrs = append(funcAddrs, fun.F)
			symFuncs[i] = fun
		}
		baselib[sym] = symFuncs
	}
}

type Func struct {
	In   []baseType
	Out  []baseType
	F    func(*Bytecode)
	addr int
}

const (
	iopIntCopy = iota
	iopStrCopy
	iopBoolCopy
)

var funcAddrs = []func(*Bytecode){
	func(p *Bytecode) { // iopIntCopy
		p.Ints[p.OpAddrs[p.pos+2]] = p.Ints[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
	func(p *Bytecode) { // iopStrCopy
		p.Strs[p.OpAddrs[p.pos+2]] = p.Strs[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
	func(p *Bytecode) { // iopBoolCopy
		p.Bools[p.OpAddrs[p.pos+2]] = p.Bools[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
}

var baselib = map[string][]Func{
	"+": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Int},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] + p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
		{
			In:  []baseType{Str, Str},
			Out: []baseType{Str},
			F: func(p *Bytecode) {
				p.Strs[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] + p.Strs[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"*": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Int},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] * p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"%": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Int},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] % p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"/": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Int},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] / p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"-": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Int},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] - p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	">": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] > p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"<": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] < p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	">=": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] >= p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"<=": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] <= p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"==": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] == p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
		{
			In:  []baseType{Str, Str},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] == p.Strs[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"!=": {
		{
			In:  []baseType{Int, Int},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] != p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
		{
			In:  []baseType{Str, Str},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] != p.Strs[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"&&": {
		{
			In:  []baseType{Bool, Bool},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] && p.Bools[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"||": {
		{
			In:  []baseType{Bool, Bool},
			Out: []baseType{Bool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] || p.Bools[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"Print": {
		{
			In: []baseType{Str},
			F: func(p *Bytecode) {
				log.Println(p.Strs[p.OpAddrs[p.pos+1]])
				p.pos += 2
			},
		},
		{
			In: []baseType{Int},
			F: func(p *Bytecode) {
				log.Println(p.Ints[p.OpAddrs[p.pos+1]])
				p.pos += 2
			},
		},
		{
			In: []baseType{Bool},
			F: func(p *Bytecode) {
				log.Println(p.Bools[p.OpAddrs[p.pos+1]])
				p.pos += 2
			},
		},
	},
	"If": {
		{
			In: []baseType{Bool},
			F: func(p *Bytecode) {
				if p.Bools[p.OpAddrs[p.pos+1]] {
					p.pos += 3
				} else {
					p.pos = p.OpAddrs[p.pos+2]
				}
			},
		},
	},
	"Goto": {
		{
			In: []baseType{Addr},
			F: func(p *Bytecode) {
				p.pos = p.Ints[p.OpAddrs[p.pos+1]]
			},
		},
	},
}
