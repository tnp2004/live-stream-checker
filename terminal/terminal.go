package terminal

import (
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tnp2004/live-stream-checker/checker"
	"github.com/tnp2004/live-stream-checker/config"
	"github.com/tnp2004/live-stream-checker/filereader"
	"github.com/tnp2004/live-stream-checker/models"
)

const (
	ERROR_STATUS    string = "error"
	LIVE_STATUS     string = "live"
	OFFLINE_STATUS  string = "offline"
	CHECKING_STATUS string = "checking"
)

const LOG_FILE_NAME = "debug.log"

type tickMsg time.Time

var once sync.Once

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type terminalModel struct {
	config      *config.Config
	channelList []*models.Channel
	width       int
	height      int
	table       table.Model
	selected    int
}

func (m terminalModel) Init() tea.Cmd {
	once.Do(func() { go m.fetchLiveStatus() })

	return tickCmd()
}

func (m terminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tickMsg:
		return m.UpdateTable()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlR:
			m.fetchLiveStatus()
			return m, tickCmd()
		case tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			}
		case tea.KeyDown:
			if m.selected < len(m.channelList)-1 {
				m.selected++
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m terminalModel) View() string {
	return DefaultViewStyles().Render(m.table.View()) + m.helpView()
}

func Run() error {
	file, err := tea.LogToFile(LOG_FILE_NAME, "debug")
	if err != nil {
		return err
	}
	defer file.Close()

	channelList := filereader.ReadChannelList()
	config := config.LoadConfig()
	m := terminalModel{
		channelList: channelList,
		table:       NewTable(channelList, 0),
		config:      config,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m terminalModel) helpView() string {
	return "\n ↑/↓: Navigate \n Ctrl+c,Esc: Quit \n Ctrl+r: recheck"
}

func (m terminalModel) fetchLiveStatus() {
	for _, ch := range m.channelList {
		ch.Status = CHECKING_STATUS
		checker := checker.New(ch, m.config)
		go func(ch *models.Channel) {
			liveStatus, err := checker.IsLive(ch.Link)
			if err != nil {
				ch.Status = ERROR_STATUS
				return
			}

			if liveStatus {
				ch.Status = LIVE_STATUS
			} else {
				ch.Status = OFFLINE_STATUS
			}
		}(ch)
	}
}
