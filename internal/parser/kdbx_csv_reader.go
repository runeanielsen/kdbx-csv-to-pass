package parser

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

type FileReader interface {
	ReadLines(path string) []string
}

type KdbxCsvReader struct{}

func (KdbxCsvReader) ReadLines(path string) []string {
	lines := make([]string, 0)

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		matched, err := regexp.MatchString(`(?:Recycle)|(?:Backup).*`, line)

		if err != nil {
			log.Fatal(err)
		}

		if !matched {
			lines = append(lines, line)
		}
	}

	return lines
}
