package parser

import (
	"fmt"
	"strings"
)

type passwordEntry struct {
	name     string
	password string
}

func Parse(filePath string, cmdExec CmdExec, fileReader FileReader) {
	csvKdbxFile := fileReader.ReadLines(filePath)

	var cmd string
	for _, line := range csvKdbxFile {
		entry := parseLine(line)
		if entry.name != "" {
			cmd += createPassCommand(entry)
		}
	}

	cmdExec.ExecuteCmd(cmd)
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
