package list

import (
	"com.steven/main/nyaa"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Items    []nyaa.Item
	focus    bool
	Max      int
	Width    int
	Cursor   int
	Selected map[int]struct{}
}

func (m Model) Focused() bool {
	return m.focus
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.focus {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			// The "up" and "k" keys move the cursor up
			case "up", "k":
				if m.Cursor > 0 {
					m.Cursor--
				}

			// The "down" and "j" keys move the cursor down
			case "down", "j":
				if m.Cursor < len(m.Items)-1 {
					m.Cursor++
				}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case " ":
				_, ok := m.Selected[m.Cursor]
				if ok {
					delete(m.Selected, m.Cursor)
				} else {
					m.Selected[m.Cursor] = struct{}{}
				}
			}
		}
	}

	return m, nil
}

func (m *Model) Search(query string) {
	m.Items = nyaa.Query(query, "seeders", "desc")
}

func (m Model) View() string {
	view := ""

	for i, item := range m.Items {
		if i == m.Max {
			break
		}

		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		name := item.Name
		if len(name) > m.Width {
			name = item.Name[0:m.Width]
		}

		view += fmt.Sprintf("\n%s[%s] %s %s", cursor, checked, name, item.Size)
	}

	return view + "\n"
}
