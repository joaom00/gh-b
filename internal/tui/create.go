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
	ti.Placeholder = "Try feature..."

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
	title := m.style.Title.MarginLeft(2).Render("Type name of the new branch")
	textInput := lipgloss.NewStyle().MarginLeft(3).Render(m.create.textinput.View())
	help := lipgloss.NewStyle().MarginLeft(3).Render(m.create.help.View(m.keyMap))

	if m.create.showConfirmInput {
		branch := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(m.create.textinput.Value())

		confirmInput := lipgloss.NewStyle().
			MarginLeft(3).
			Render(lipgloss.JoinHorizontal(
				lipgloss.Left,
				fmt.Sprintf("Create new branch \"%s\"? [y/n]:", branch),
				m.create.confirmInput.View(),
			))

		return lipgloss.NewStyle().
			MarginTop(1).
			Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", textInput, "\n", confirmInput, "\n", help))
	}

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", textInput, "\n", help))
}
