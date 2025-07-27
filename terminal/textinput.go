package terminal

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type AddInput struct {
	Index int
	Input []textinput.Model
}

func (ai *AddInput) ClearInput() {
	for i := range ai.Input {
		ai.Input[i].SetValue("")
	}
}

func NewAddInput() AddInput {
	return AddInput{
		Index: 0,
		Input: []textinput.Model{
			newTextInput("enter channel name"),
			newTextInput("place your channel url"),
		},
	}
}

func newTextInput(placeHolder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeHolder
	ti.Width = 50
	return ti
}

func (m terminalModel) AddHelpView() string {
	return "\n ↑/↓: Navigate \n Enter: add to channel list \n Esc: back to checker \n Ctrl+c: Quit"
}
