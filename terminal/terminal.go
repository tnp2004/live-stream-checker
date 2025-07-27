package terminal

import (
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tnp2004/live-stream-checker/checker"
	"github.com/tnp2004/live-stream-checker/config"
	"github.com/tnp2004/live-stream-checker/filereader"
	"github.com/tnp2004/live-stream-checker/models"
)

const (
	CHECK_STAGE int = iota
	ADD_STAGE
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
	stage       int
	addInput    AddInput
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
		case tea.KeyEnter:
			_ = filereader.AddChannel(m.addInput.Input[0].Value(), m.addInput.Input[1].Value())
			m.stage = CHECK_STAGE
			channelList := filereader.ReadChannelList()
			m.channelList = append(m.channelList, channelList[len(channelList)-1])
			m.fetchLiveStatus()
			m.addInput.ClearInput()
			return m, tickCmd()
		case tea.KeyCtrlA:
			m.stage = ADD_STAGE
			m.addInput.Index = 0
			m.addInput.Input[m.addInput.Index], cmd = m.addInput.Input[m.addInput.Index].Update(msg)
			return m, cmd
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			switch m.stage {
			case CHECK_STAGE:
				return m, tea.Quit
			case ADD_STAGE:
				m.stage = CHECK_STAGE
				return m, nil
			}
		case tea.KeyCtrlR:
			m.fetchLiveStatus()
			return m, tickCmd()
		case tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			}
			if m.addInput.Index > 0 {
				m.addInput.Index--
			}
		case tea.KeyDown:
			if m.selected < len(m.channelList)-1 {
				m.selected++
			}
			if m.addInput.Index < len(m.addInput.Input)-1 {
				m.addInput.Index++
			}
		}
	}

	switch m.stage {
	case CHECK_STAGE:
		m.table, cmd = m.table.Update(msg)
	case ADD_STAGE:
		for i := range m.addInput.Input {
			if m.addInput.Index == i {
				m.addInput.Input[i].Focus()
				m.addInput.Input[i].Cursor.SetMode(cursor.CursorBlink)
				continue
			}
			m.addInput.Input[i].Blur()
			m.addInput.Input[i].Cursor.SetMode(cursor.CursorHide)
		}
		m.addInput.Input[m.addInput.Index], cmd = m.addInput.Input[m.addInput.Index].Update(msg)
		return m, cmd
	}
	return m, cmd
}

func (m terminalModel) View() string {
	var s string
	switch m.stage {
	case CHECK_STAGE:
		s = DefaultViewStyles().Render(m.table.View()) + m.checkHelpView()
	case ADD_STAGE:
		var b strings.Builder

		for i := range m.addInput.Input {
			b.WriteString(m.addInput.Input[i].View())
			if i < len(m.addInput.Input)-1 {
				b.WriteRune('\n')
			}
		}

		s = DefaultViewStyles().Render(b.String()) + m.AddHelpView()
	}

	return s
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
		stage:       CHECK_STAGE,
		channelList: channelList,
		table:       NewTable(channelList, 0),
		addInput:    NewAddInput(),
		config:      config,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m terminalModel) checkHelpView() string {
	return "\n ↑/↓: Navigate \n Ctrl+a: Add channel \n Ctrl+c,Esc: Quit \n Ctrl+r: recheck"
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
