package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type passwordEntry struct {
	name     string
	password string
}

func main() {
	filePath := "/home/notation/kdbx/pass"

	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	csvKdbxFile := readFile(filePath)

	var wg sync.WaitGroup

	for _, line := range csvKdbxFile {
		wg.Add(1)
		// This is ugly but needed to make sure that the pass command is not locked
		// Any suggestion please say so
		time.Sleep(2 * time.Millisecond)
		go handleKdbxLine(&wg, line)
	}

	wg.Wait()
}

func handleKdbxLine(wg *sync.WaitGroup, line string) {
	defer wg.Done()
	entry := parseLine(line)

	if entry.name != "" {
		fmt.Printf("Inserting pass %s\n", entry.name)
		store(entry)
	}
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
	return strings.ReplaceAll(text, "'", "\\'")
}

func removeQuotes(text string) string {
	return strings.ReplaceAll(text, "\"", "")
}

func cleanName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func store(entry passwordEntry) {
	cmd := fmt.Sprintf("echo '%s' | pass insert '%s' -e", entry.password, entry.name)
	_, err := exec.Command("fish", "-c", cmd).Output()

	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", cmd)
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
