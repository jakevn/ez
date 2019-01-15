package ez

type Bytecode struct {
	OpAddrs []int    `json:"op_addrs,omitempty"`
	Ints    []int    `json:"ints,omitempty"`
	Strs    []string `json:"strs,omitempty"`
	Bools   []bool   `json:"bools,omitempty"`
	pos     int
}
