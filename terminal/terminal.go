package terminal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tnp2004/live-stream-checker/filereader"
	"github.com/tnp2004/live-stream-checker/models"
)

const LOG_FILE_NAME = "debug.log"

type terminalModel struct {
	channelList []*models.Channel
	width       int
	height      int
}

func (m terminalModel) Init() tea.Cmd {
	return nil
}

func (m terminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m terminalModel) View() string {
	return "hello"
}

func Run() error {
	file, err := tea.LogToFile(LOG_FILE_NAME, "debug")
	if err != nil {
		return err
	}
	defer file.Close()

	channelList := filereader.ReadChannelList()
	m := terminalModel{channelList: channelList}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
