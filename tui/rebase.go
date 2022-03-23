package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaom00/gh-b/git"
)

type rebaseModel struct {
	help         help.Model
	confirmInput textinput.Model
}

func newRebaseModel() *rebaseModel {
	ci := textinput.New()
	ci.CharLimit = 1

	return &rebaseModel{
		help:         help.New(),
		confirmInput: ci,
	}
}

func rebaseUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.rebase.confirmInput.Value() {
			case "y", "Y", "":
				if i, ok := m.list.SelectedItem().(item); ok {
					i.Name = strings.TrimSuffix(i.Name, "*")
					out := git.RebaseBranch(i.Name)

					fmt.Println(m.styles.NormalTitle.Copy().MarginTop(1).Render(out))
					return m, tea.Quit
				}

			case "n", "N":
				m.rebase.confirmInput.Reset()
				m.state = browsing
				m.updateKeybindins()

			default:
				m.rebase.confirmInput.SetValue("")
			}
		case key.Matches(msg, m.keyMap.Cancel):
			m.rebase.confirmInput.Reset()
			m.state = browsing
			m.updateKeybindins()
		}
	case tea.WindowSizeMsg:
		m.rebase.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.rebase.confirmInput, cmd = m.rebase.confirmInput.Update(msg)

	return m, cmd
}

func (m Model) rebaseView() string {
	var branchName string

	i, ok := m.list.SelectedItem().(item)
	if ok {
		branchName = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(i.Name)
	}

	label := fmt.Sprintf("Do you really wanna rebase branch \"%s\"? [Y/n]", branchName)

	confirmInput := lipgloss.JoinHorizontal(
		lipgloss.Left,
		label,
		m.rebase.confirmInput.View(),
	)

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(4).
		Render(lipgloss.JoinVertical(lipgloss.Left, confirmInput, "\n", m.rebase.help.View(m.keyMap)))
}
