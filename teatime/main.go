package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	_WEB_SOCKET_URL = "ws://localhost:8080/ws"
)

// flags for the login
var (
	USERNAME  = flag.String("username", "", "username")
	PASSWORD  = flag.String("password", "", "password")
	ROOM_NAME = flag.String("room", "", "room name")
	ROOM_PASS = flag.String("roompass", "", "room password")
)

func main() {
	flag.Parse()
	if *USERNAME == "" || *PASSWORD == "" || *ROOM_NAME == "" || *ROOM_PASS == "" {
		log.Fatal("Please provide all the flags")
	}
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	err           error
	socket_client *client
}

func initialModel() *model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(40)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = true

	vp := viewport.New(40, 10)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(true)

	m := model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
	m.socket_client = InitOrDie(m.AddSocketMessage)
	return &m
}

func (m *model) AddSocketMessage(msg []byte) {
	payload := map[string]interface{}{}
	err := json.Unmarshal(msg, &payload)
	if err != nil {
		return
	}
	sender, ok := payload["user"].(string)
	if !ok {
		return
	}
	message, ok := payload["message"].(string)
	if !ok {
		return
	}
	m.messages = append(m.messages, m.senderStyle.Render(sender+": ")+message)
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()
}

func (m *model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			sendObj := map[string]string{
				"user":    *USERNAME,
				"message": m.textarea.Value(),
				"room":    *ROOM_NAME,
			}
			m.socket_client.SendMessageJSON(sendObj)
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
