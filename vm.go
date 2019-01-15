package ez

import "log"

func Run(p *Bytecode) {
	for p.pos < len(p.OpAddrs) {
		switch p.OpAddrs[p.pos] {
		case 0: // 0: iopIntCopy (int int)
			p.Ints[p.OpAddrs[p.pos+2]] = p.Ints[p.OpAddrs[p.pos+1]]
			p.pos += 3
		case 1: // 1: iopStrCopy (str str)
			p.Strs[p.OpAddrs[p.pos+2]] = p.Strs[p.OpAddrs[p.pos+1]]
			p.pos += 3
		case 2: // 2: iopBoolCopy (bool bool)
			p.Bools[p.OpAddrs[p.pos+2]] = p.Bools[p.OpAddrs[p.pos+1]]
			p.pos += 3
		case 3: // 3: != (int int) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] != p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 4: // 4: != (str str) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] != p.Strs[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 5: // 5: % (int int) -> int
			p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] % p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 6: // 6: && (bool bool) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] && p.Bools[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 7: // 7: * (int int) -> int
			p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] * p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 8: // 8: + (int int) -> int
			p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] + p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 9: // 9: + (str str) -> str
			p.Strs[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] + p.Strs[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 10: // 10: - (int int) -> int
			p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] - p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 11: // 11: / (int int) -> int
			p.Ints[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] / p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 12: // 12: < (int int) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] < p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 13: // 13: <= (int int) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] <= p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 14: // 14: == (int int) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] == p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 15: // 15: == (str str) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Strs[p.OpAddrs[p.pos+1]] == p.Strs[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 16: // 16: > (int int) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] > p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 17: // 17: >= (int int) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Ints[p.OpAddrs[p.pos+1]] >= p.Ints[p.OpAddrs[p.pos+2]]
			p.pos += 4
		case 18: // 18: goto (addr)
			p.pos = p.Ints[p.OpAddrs[p.pos+1]]
		case 19: // 19: if (bool)
			if p.Bools[p.OpAddrs[p.pos+1]] {
				p.pos += 3
			} else {
				p.pos = p.OpAddrs[p.pos+2]
			}
		case 20: // 20: print (str)
			log.Println(p.Strs[p.OpAddrs[p.pos+1]])
			p.pos += 2
		case 21: // 21: print (int)
			log.Println(p.Ints[p.OpAddrs[p.pos+1]])
			p.pos += 2
		case 22: // 22: print (bool)
			log.Println(p.Bools[p.OpAddrs[p.pos+1]])
			p.pos += 2
		case 23: // 23: || (bool bool) -> bool
			p.Bools[p.OpAddrs[p.pos+3]] = p.Bools[p.OpAddrs[p.pos+1]] || p.Bools[p.OpAddrs[p.pos+2]]
			p.pos += 4
		}
	}
}
