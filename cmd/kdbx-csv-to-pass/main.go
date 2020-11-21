package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/runeanielsen/kdbx-csv-to-pass/internal/parser"
)

func main() {
	file := flag.String("flag", "", "Path to exported kdbx.csv")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Printf("Starting import of passwords from file %s\n", *file)
	parser.Parse(*file, parser.CmdExecBash{}, parser.KdbxCsvReader{})
	fmt.Println("Finished import passwords")
}
