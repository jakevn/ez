package ez

import (
	"log"
	"sort"
)

func init() {
	syms := make([]string, len(baselib))
	for sym, _ := range baselib {
		syms = append(syms, sym)
	}
	sort.Strings(syms)
	for _, sym := range syms {
		symFuncs := baselib[sym]
		for i, fun := range symFuncs {
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
	func(p *Bytecode) {
		p.Ints[p.OpAddrs[p.pos+2]] = p.Ints[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
	func(p *Bytecode) {
		p.Strs[p.OpAddrs[p.pos+2]] = p.Strs[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
	func(p *Bytecode) {
		p.Bools[p.OpAddrs[p.pos+2]] = p.Bools[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
}

var baselib = map[string][]Func{
	"+": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btInt},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] + p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
		{
			In:  []baseType{btStr, btStr},
			Out: []baseType{btStr},
			F: func(p *Bytecode) {
				p.Strs[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] + p.Strs[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"*": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btInt},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] * p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"%": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btInt},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] % p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"/": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btInt},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] / p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"-": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btInt},
			F: func(p *Bytecode) {
				p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] - p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	">": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] > p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"<": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] < p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	">=": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] >= p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"<=": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] <= p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"==": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] == p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
		{
			In:  []baseType{btStr, btStr},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] == p.Strs[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"!=": {
		{
			In:  []baseType{btInt, btInt},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] != p.Ints[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
		{
			In:  []baseType{btStr, btStr},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] != p.Strs[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"&&": {
		{
			In:  []baseType{btBool, btBool},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] && p.Bools[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"||": {
		{
			In:  []baseType{btBool, btBool},
			Out: []baseType{btBool},
			F: func(p *Bytecode) {
				p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] || p.Bools[p.OpAddrs[p.pos+2]]
				p.pos += 4
			},
		},
	},
	"Print": {
		{
			In: []baseType{btStr},
			F: func(p *Bytecode) {
				log.Println(p.Strs[p.OpAddrs[p.pos+1]])
				p.pos += 2
			},
		},
		{
			In: []baseType{btInt},
			F: func(p *Bytecode) {
				log.Println(p.Ints[p.OpAddrs[p.pos+1]])
				p.pos += 2
			},
		},
		{
			In: []baseType{btBool},
			F: func(p *Bytecode) {
				log.Println(p.Bools[p.OpAddrs[p.pos+1]])
				p.pos += 2
			},
		},
	},
}
