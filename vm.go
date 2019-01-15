package ez

import "log"

func Run(p *Bytecode) {
	for p.pos < len(p.OpAddrs) {
		funcAddrs[p.OpAddrs[p.pos]](p)
	}
}

var funcAddrs = []func(*Bytecode){
	func(p *Bytecode) { // 0: iopIntCopy (int int)
		p.Ints[p.OpAddrs[p.pos+2]] = p.Ints[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
	func(p *Bytecode) { // 1: iopStrCopy (str str)
		p.Strs[p.OpAddrs[p.pos+2]] = p.Strs[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
	func(p *Bytecode) { // 2: iopBoolCopy (bool bool)
		p.Bools[p.OpAddrs[p.pos+2]] = p.Bools[p.OpAddrs[p.pos+1]]
		p.pos += 3
	},
	func(p *Bytecode) { // 3: != (int int) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] != p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 4: != (str str) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] != p.Strs[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 5: % (int int) -> int
		p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] % p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 6: && (bool bool) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] && p.Bools[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 7: * (int int) -> int
		p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] * p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 8: + (int int) -> int
		p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] + p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 9: + (str str) -> str
		p.Strs[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] + p.Strs[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 10: - (int int) -> int
		p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] - p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 11: / (int int) -> int
		p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] / p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 12: < (int int) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] < p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 13: <= (int int) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] <= p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 14: == (int int) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] == p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 15: == (str str) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] == p.Strs[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 16: > (int int) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] > p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 17: >= (int int) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] >= p.Ints[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
	func(p *Bytecode) { // 18: goto (addr)
		p.pos = p.Ints[p.OpAddrs[p.pos+1]]
	},
	func(p *Bytecode) { // 19: if (bool)
		if p.Bools[p.OpAddrs[p.pos+1]] {
			p.pos += 3
		} else {
			p.pos = p.OpAddrs[p.pos+2]
		}
	},
	func(p *Bytecode) { // 20: print (str)
		log.Println(p.Strs[p.OpAddrs[p.pos+1]])
		p.pos += 2
	},
	func(p *Bytecode) { // 21: print (int)
		log.Println(p.Ints[p.OpAddrs[p.pos+1]])
		p.pos += 2
	},
	func(p *Bytecode) { // 22: print (bool)
		log.Println(p.Bools[p.OpAddrs[p.pos+1]])
		p.pos += 2
	},
	func(p *Bytecode) { // 23: || (bool bool) -> bool
		p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] || p.Bools[p.OpAddrs[p.pos+2]]
		p.pos += 4
	},
}
