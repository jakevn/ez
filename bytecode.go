package ez

type Bytecode struct {
	OpAddrs []int
	Ints    []int
	Strs    []string
	Bools   []bool
	pos     int
}
