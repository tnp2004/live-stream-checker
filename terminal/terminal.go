package terminal

import (
	tea "github.com/charmbracelet/bubbletea"
)

const LOG_FILE_NAME = "debug.log"

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m model) View() string {
	return ""
}

func Run() error {
	file, err := tea.LogToFile(LOG_FILE_NAME, "debug")
	if err != nil {
		return err
	}
	defer file.Close()
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
