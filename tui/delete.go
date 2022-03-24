package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaom00/gh-b/git"
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
				if i, ok := m.list.SelectedItem().(item); ok {
					i.Name = strings.TrimSuffix(i.Name, "*")
					out := git.DeleteBranch(i.Name)

					fmt.Println(m.styles.NormalTitle.Render(out))

					time.Sleep(time.Second * 1)

					m.updateListItem()
					m.state = browsing
					m.keyMap.State = "browsing"
					m.updateKeybindins()
				}

			case "n", "N":
				m.delete.confirmInput.Reset()
				m.state = browsing
				m.keyMap.State = "browsing"
				m.updateKeybindins()

			default:
				m.delete.confirmInput.SetValue("")
			}

		case key.Matches(msg, m.keyMap.Cancel):
			m.delete.confirmInput.Reset()
			m.state = browsing
			m.keyMap.State = "browsing"
			m.updateKeybindins()
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

	if i, ok := m.list.SelectedItem().(item); ok {
		branchName = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(i.Name)
	}

	label := fmt.Sprintf("Do you really wanna delete branch \"%s\"? [Y/n]", branchName)

	confirmInput := lipgloss.JoinHorizontal(
		lipgloss.Left,
		label,
		m.delete.confirmInput.View(),
	)

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(4).
		Render(lipgloss.JoinVertical(lipgloss.Left, confirmInput, "\n", m.delete.help.View(m.keyMap)))
}
