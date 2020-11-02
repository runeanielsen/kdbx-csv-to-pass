package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type passwordEntry struct {
	name     string
	password string
}

func Parse(filePath string) {
	csvKdbxFile := readFile(filePath)

	var cmd string

	for _, line := range csvKdbxFile {
		entry := parseLine(line)
		if entry.name != "" {
			cmd += createPassCommand(entry)
		}
	}

	err := exec.Command("bash", "-c", cmd).Run()

	if err != nil {
		log.Fatalf("Failed to execute command: %s - %s\n", cmd, err)
	}
}

func parseLine(line string) passwordEntry {
	splitted := strings.Split(line, ",")

	if len(splitted) > 4 {
		entryName, username, password := splitted[1], splitted[2], splitted[3]

		entryName = removeQuotes(entryName)
		username = removeQuotes(username)
		password = removeQuotes(password)
		password = escapeSingleQuotes(password)

		newEntryName := entryName
		if username != "" {
			newEntryName += "-" + username
		}

		newEntryName = standardize(newEntryName)

		return passwordEntry{
			name:     newEntryName,
			password: password,
		}
	}

	return passwordEntry{}
}

func escapeSingleQuotes(text string) string {
	return strings.ReplaceAll(text, "'", "''")
}

func removeQuotes(text string) string {
	return strings.ReplaceAll(text, "\"", "")
}

func standardize(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "_")
	return text
}

func createPassCommand(entry passwordEntry) string {
	return fmt.Sprintf("echo '%s' | pass insert '%s' -e;", entry.password, entry.name)
}

func readFile(filepath string) []string {
	lines := make([]string, 0)

	file, err := os.Open(filepath)
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
