package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaom00/gh-b/internal/git"
)

type deleteModel struct {
	help         help.Model
	confirmInput textinput.Model
}

func newDeleteModel() *deleteModel {
	ci := textinput.New()
	ci.CharLimit = 1

	return &deleteModel{
		help:         help.New(),
		confirmInput: ci,
	}
}

func deleteUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.delete.confirmInput.Value() {
			case "y", "Y", "":
				i, ok := m.list.SelectedItem().(item)
				if ok {
					out := git.DeleteBranch(string(i))

					fmt.Println("\n", out)
					return m, tea.Quit
				}

			case "n", "N":
				m.delete.confirmInput.Reset()
				m.state = browsing

			default:
				m.delete.confirmInput.SetValue("")
			}
		case key.Matches(msg, m.keyMap.Cancel):
			m.delete.confirmInput.Reset()
			m.state = browsing
		}
	case tea.WindowSizeMsg:
		m.delete.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.delete.confirmInput, cmd = m.delete.confirmInput.Update(msg)

	return m, cmd
}

func (m Model) deleteView() string {
	var branchName string

	i, ok := m.list.SelectedItem().(item)
	if ok {
		branchName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Render(string(i))
	}

	label := fmt.Sprintf("Do you really wanna delete branch \"%s\"? [Y/n]", branchName)

	confirmInput := lipgloss.JoinHorizontal(
		lipgloss.Left,
		label,
		m.delete.confirmInput.View(),
	)

	return lipgloss.NewStyle().
		Margin(1, 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, confirmInput, "\n", m.delete.help.View(m.keyMap)))
}
