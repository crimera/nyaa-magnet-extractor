package main

import (
	"fmt"
	"os"

	//"com.steven/main/nyaa"
	"com.steven/main/nyaa"
	"com.steven/main/ui/list"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list      list.Model
	textInput textinput.Model
}

func initialModel() model {
	list := list.Model{
		Selected: make(map[int]struct{}),
		Width:    100,
		Max:      20,
	}

	textinput := textinput.New()
	textinput.Focus()

	return model{
		list:      list,
		textInput: textinput,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.textInput.Blur()
			m.list.Focus()
		case "ctrl+c", "q":
			if !m.textInput.Focused() {
				return m, tea.Quit
			}
		case "a":
			if !m.textInput.Focused() {
				AddTorents(m.list)
			}
		case "enter":
			m.list.Search(m.textInput.Value())
			m.list.Focus()
			m.textInput.Blur()
		case "/":
			m.list.Blur()
			m.textInput.Focus()
		}
	}

	m.list, cmd = m.list.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func AddTorents(m list.Model) {
	var magnets []string
	for i, item := range m.Items {
		if _, ok := m.Selected[i]; ok {
			magnets = append(magnets, item.Magnet)
		}
	}

	for _, magnet := range magnets {
		nyaa.AddTorrent("/home/dingle/Videos/Anime", magnet)
	}
}

func (m model) View() string {
	return m.textInput.View() + m.list.View()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
