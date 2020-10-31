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

	store(cmd)
}

func parseLine(line string) passwordEntry {
	splitted := strings.Split(line, ",")

	if len(splitted) > 4 {
		name, username, password := splitted[1], splitted[2], splitted[3]

		name = removeQuotes(name)
		username = removeQuotes(username)
		password = removeQuotes(password)
		password = escapeSingleQuotes(password)

		newName := name
		if username != "" {
			newName += "-" + username
		}

		newName = cleanName(newName)

		return passwordEntry{
			name:     newName,
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

func cleanName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func createPassCommand(entry passwordEntry) string {
	return fmt.Sprintf("echo '%s' | pass insert '%s' -e;", entry.password, entry.name)
}

func store(cmd string) {
	err := exec.Command("bash", "-c", cmd).Run()

	if err != nil {
		log.Fatalf("Failed to execute command: %s - %s\n", cmd, err)
	}
}

func readFile(filepath string) []string {
	lines := make([]string, 0)

	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

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
