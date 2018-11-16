package ez

func Run(p *Bytecode) {
	for p.pos < len(p.OpAddrs) {
		funcAddrs[p.OpAddrs[p.pos]](p)
	}
}
