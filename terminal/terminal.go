package terminal

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tnp2004/live-stream-checker/checker"
	"github.com/tnp2004/live-stream-checker/config"
	"github.com/tnp2004/live-stream-checker/filereader"
	"github.com/tnp2004/live-stream-checker/models"
)

const (
	ERROR_STATUS   string = "error"
	LIVE_STATUS    string = "live"
	OFFLINE_STATUS string = "offline"
)

const LOG_FILE_NAME = "debug.log"

type terminalModel struct {
	channelList []*models.Channel
	width       int
	height      int
	table       table.Model
}

func (m terminalModel) fetchLiveStatus() {
	cfg := config.LoadConfig()
	for _, ch := range m.channelList {
		checker := checker.New(ch, cfg)
		liveStatus, err := checker.IsLive(ch.Link)
		if err != nil {
			ch.Status = ERROR_STATUS
		}

		if liveStatus {
			ch.Status = LIVE_STATUS
		} else {
			ch.Status = OFFLINE_STATUS
		}
	}
}

func (m terminalModel) Init() tea.Cmd {
	m.fetchLiveStatus()
	return nil
}

func (m terminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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

	m.table = m.NewTable()
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m terminalModel) View() string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).Render(m.table.View() + "\n")
}

func (m terminalModel) NewTable() table.Model {
	rows := make([]table.Row, 0, len(m.channelList))
	var longestChName int
	for _, ch := range m.channelList {
		len := len(ch.Name)
		if len > longestChName {
			longestChName = len
		}
		row := table.Row{ch.Name, ch.Platform, ch.Status}
		rows = append(rows, row)
	}

	columns := []table.Column{
		{Title: "Channel", Width: longestChName},
		{Title: "Platform", Width: 8},
		{Title: "Status", Width: 7},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
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
