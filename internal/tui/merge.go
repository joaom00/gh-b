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

type mergeModel struct {
	help         help.Model
	confirmInput textinput.Model
}

func newMergeModel() *mergeModel {
	ci := textinput.New()
	ci.CharLimit = 1

	return &mergeModel{
		help:         help.New(),
		confirmInput: ci,
	}
}

func mergeUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.delete.confirmInput.Value() {
			case "y", "Y", "":
				i, ok := m.list.SelectedItem().(item)
				if ok {
					out := git.MergeBranch(i.Name)

					fmt.Println("\n", out)
					return m, tea.Quit
				}

			case "n", "N":
				m.merge.confirmInput.Reset()
				m.state = browsing

			default:
				m.merge.confirmInput.SetValue("")
			}
		case key.Matches(msg, m.keyMap.Cancel):
			m.merge.confirmInput.Reset()
			m.state = browsing
		}

	case tea.WindowSizeMsg:
		m.merge.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.merge.confirmInput, cmd = m.merge.confirmInput.Update(msg)

	return m, cmd
}

func (m Model) mergeView() string {
	var branchName string

	i, ok := m.list.SelectedItem().(item)
	if ok {
		branchName = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(i.Name)
	}

	label := fmt.Sprintf("Do you really wanna merge branch \"%s\"? [Y/n]", branchName)

	confirmInput := lipgloss.JoinHorizontal(
		lipgloss.Left,
		label,
		m.merge.confirmInput.View(),
	)

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(4).
		Render(lipgloss.JoinVertical(lipgloss.Left, confirmInput, "\n", m.merge.help.View(m.keyMap)))
}
