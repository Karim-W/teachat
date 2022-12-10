package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginModel struct {
	viewport    viewport.Model
	username    textarea.Model
	password    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialLoginModel() LoginModel {
	ta := textarea.New()
	ta.Placeholder = "Enter your login..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return LoginModel{
		username:    ta,
		password:    ta,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		}
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m LoginModel) View() string {
	return m.viewport.View()
}
