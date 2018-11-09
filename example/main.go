package main

import (
	"log"
	"time"

	"github.com/jakevn/ez"
)

const testProg = `
testInt = 3 + 9
anotherInt = 12 + testInt
aString = 'hello'
anotherString = aString + ' world'
subtraction = 100 - 23
division = 30 / 10
multiplication = 11 * 13
`

func main() {
	log.SetFlags(0)
	start := time.Now()
	prog := &ez.Prog{
		IntIDAddr:   map[string]int{},
		StrIDAddr:   map[string]int{},
		FuncIDImpls: map[string][]ez.Func{},
	}
	if err := prog.AddOps(ez.Prelude()); err != nil {
		log.Fatal(err)
	}
	if err := prog.Parse([]byte(testProg)); err != nil {
		log.Fatal(err)
	}
	println("parsing elapsed ns:", time.Since(start).Nanoseconds())
	start = time.Now()
	prog.Run()
	println("run elapsed ns:", time.Since(start).Nanoseconds())

	prog.DebugPrintSymbol("testInt")
	prog.DebugPrintSymbol("anotherInt")
	prog.DebugPrintSymbol("anotherString")
	prog.DebugPrintSymbol("subtraction")
	prog.DebugPrintSymbol("division")
	prog.DebugPrintSymbol("multiplication")
}
