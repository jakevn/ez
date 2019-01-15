package main

import (
	"io"
	"log"
	"os"
	"strings"

	"encoding/json"

	"github.com/jakevn/ez"
)

func main() {
	log.SetFlags(0)
	filePath := os.Args[1]
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var bc ez.Bytecode
	if strings.HasSuffix(filePath, ".ezc") {
		bc, err = decode(file)
	} else {
		bc, err = ez.Parse(file)
		if err == nil && hasFlag("c") {
			if err := saveByteCode(bc); err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	ez.Run(&bc)
}

func saveByteCode(bc ez.Bytecode) error {
	filePath := os.Args[1]
	filePath = strings.Replace(filePath, ".ez", "", -1)
	outFilePath := filePath + ".ezc"
	file, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(bc)
}

func decode(reader io.Reader) (ez.Bytecode, error) {
	var ezb ez.Bytecode
	return ezb, json.NewDecoder(reader).Decode(&ezb)
}

func hasFlag(flag string) bool {
	_, ok := getFlagParam(flag)
	return ok
}

func getFlagParam(flag string) (string, bool) {
	dashFlag := "-" + flag
	doubleDashFlag := "--" + flag
	var param string
	var hasFlag bool
	for _, arg := range os.Args {
		if hasFlag {
			if arg[0] == '-' {
				break
			}
			if param == "" {
				param = arg
			} else {
				param += " " + arg
			}
		} else {
			splitArg := strings.Split(arg, "=")
			if splitArg[0] != dashFlag && splitArg[0] != doubleDashFlag {
				continue
			}
			hasFlag = true
		}
	}
	return param, hasFlag
}
