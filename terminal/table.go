package terminal

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tnp2004/live-stream-checker/models"
)

func NewTable(channelList []*models.Channel, selected int) table.Model {
	rows := make([]table.Row, 0, len(channelList))
	var longestChName int
	for _, ch := range channelList {
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
		{Title: "Status", Width: 11},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	t.SetCursor(selected)

	s := DefaultTableStyles()
	t.SetStyles(s)

	return t
}

func (m terminalModel) CheckView() string {
	return DefaultViewStyles().Render(m.table.View()) + m.checkHelpView()
}

func (m terminalModel) UpdateTable() (tea.Model, tea.Cmd) {
	isAllChStatusUpdated := true
	m.table = NewTable(m.channelList, m.selected)
	for _, ch := range m.channelList {
		if ch.Status == CHECKING_STATUS {
			isAllChStatusUpdated = false
		}
	}
	var cmd tea.Cmd
	if isAllChStatusUpdated {
		cmd = nil
	} else {
		cmd = tickCmd()
	}

	return m, cmd
}
