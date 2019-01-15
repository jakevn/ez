package ez

type Func struct {
	In   []baseType
	Out  []baseType
	addr int
}

const (
	iopIntCopy = iota
	iopStrCopy
	iopBoolCopy
)

var baselib = map[string][]Func{
	"!=": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Bool},
			addr: 3,
		},
		{
			In:   []baseType{Str, Str},
			Out:  []baseType{Bool},
			addr: 4,
		},
	},
	"%": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Int},
			addr: 5,
		},
	},
	"&&": {
		{
			In:   []baseType{Bool, Bool},
			Out:  []baseType{Bool},
			addr: 6,
		},
	},
	"*": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Int},
			addr: 7,
		},
	},
	"+": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Int},
			addr: 8,
		},
		{
			In:   []baseType{Str, Str},
			Out:  []baseType{Str},
			addr: 9,
		},
	},
	"-": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Int},
			addr: 10,
		},
	},
	"/": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Int},
			addr: 11,
		},
	},
	"<": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Bool},
			addr: 12,
		},
	},
	"<=": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Bool},
			addr: 13,
		},
	},
	"==": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Bool},
			addr: 14,
		},
		{
			In:   []baseType{Str, Str},
			Out:  []baseType{Bool},
			addr: 15,
		},
	},
	">": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Bool},
			addr: 16,
		},
	},
	">=": {
		{
			In:   []baseType{Int, Int},
			Out:  []baseType{Bool},
			addr: 17,
		},
	},
	"goto": {
		{
			In:   []baseType{Addr},
			addr: 18,
		},
	},
	"if": {
		{
			In:   []baseType{Bool},
			addr: 19,
		},
	},
	"print": {
		{
			In:   []baseType{Str},
			addr: 20,
		},
		{
			In:   []baseType{Int},
			addr: 21,
		},
		{
			In:   []baseType{Bool},
			addr: 22,
		},
	},
	"||": {
		{
			In:   []baseType{Bool, Bool},
			Out:  []baseType{Bool},
			addr: 23,
		},
	},
}
