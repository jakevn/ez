package main

import (
	"bytes"
	"log"

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
modulo = 9 % 7
intLessThan = division < multiplication
`

func main() {
	bc, err := ez.Parse(bytes.NewReader([]byte(testProg)))
	if err != nil {
		log.Fatal(err)
	}
	ez.Run(&bc)
}
