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

type createModel struct {
	help             help.Model
	textinput        textinput.Model
	confirmInput     textinput.Model
	showConfirmInput bool
}

func newCreateModel() *createModel {
	ti := textinput.New()

	ci := textinput.New()
	ci.CharLimit = 1

	return &createModel{
		help:             help.New(),
		textinput:        ti,
		confirmInput:     ci,
		showConfirmInput: false,
	}
}

func createUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.create.confirmInput.Value() {
			case "y", "Y":
				out := git.CreateBranch(m.create.textinput.Value())

				fmt.Println("\n", out)
				return m, tea.Quit
			case "n", "N":
				m.create.textinput.Reset()
				m.create.confirmInput.Reset()
				m.create.showConfirmInput = false
				m.state = browsing
				return m, nil
			default:
				m.create.confirmInput.SetValue("")
			}

			m.create.textinput.Blur()
			m.create.showConfirmInput = true
			m.create.confirmInput.Focus()

		case key.Matches(msg, m.keyMap.Cancel):
			m.create.textinput.Reset()
			m.create.confirmInput.Reset()
			m.create.showConfirmInput = false
			m.state = browsing
		}
	case tea.WindowSizeMsg:
		m.create.help.Width = msg.Width
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.create.textinput, cmd = m.create.textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.create.confirmInput, cmd = m.create.confirmInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) createView() string {
	label := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		Render("Type name of the new branch:")

	if m.create.showConfirmInput {
		confirmInput := lipgloss.JoinHorizontal(
			lipgloss.Left,
			fmt.Sprintf("Create new branch \"%s\"? [y/n]:", m.create.textinput.Value()),
			m.create.confirmInput.View(),
		)
		return lipgloss.NewStyle().
			Margin(2, 3).
			Render(lipgloss.JoinVertical(lipgloss.Left, label, "\n", m.create.textinput.View(), "\n", confirmInput, "\n", m.create.help.View(m.keyMap)))
	}

	return lipgloss.NewStyle().
		Margin(2, 3).
		Render(lipgloss.JoinVertical(lipgloss.Left, label, "\n", m.create.textinput.View(), "\n", m.create.help.View(m.keyMap)))
}
