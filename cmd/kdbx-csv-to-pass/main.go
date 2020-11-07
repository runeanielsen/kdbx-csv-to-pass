package main

import (
	"fmt"
	"log"
	"os"

	"github.com/runeanielsen/kdbx-csv-to-pass/internal/parser"
)

func main() {
	if len(os.Args) > 1 {
		var filePath string = os.Args[1]
		fmt.Printf("Starting import of passwords from file %s\n", filePath)
		parser.Parse(filePath, parser.CmdExecBash{}, parser.KdbxCsvReader{})
	} else {
		log.Fatal("Parameter for filepath is missing")
	}

	fmt.Println("Finished import passwords")
}
