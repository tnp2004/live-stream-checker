package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const LOG_FILE_NAME = "debug.log"

func main() {
	file, err := tea.LogToFile(LOG_FILE_NAME, "debug")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer file.Close()
}
